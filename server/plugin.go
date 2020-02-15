package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"image/jpeg"
	_ "image/jpeg"

	facesLib "github.com/juanfran/mattermost-img-face-swap/server/faces"
	"github.com/juanfran/mattermost-img-face-swap/server/swap"
	"github.com/juanfran/mattermost-img-face-swap/server/utils"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/segmentio/ksuid"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	router *mux.Router
}

var (
	maxInMemoryMemes                         = 10
	generatedMemes    map[string]image.Image = make(map[string]image.Image)
	generatedMemesIds []string
	currentImagePath  string
)

func generateID() string {
	return ksuid.New().String()
}

func serveImg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=604800")

	vars := mux.Vars(r)
	memeID := vars["name"]

	if err := jpeg.Encode(w, generatedMemes[memeID], &jpeg.Options{
		Quality: 90,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) generateMeme(faces []*facesLib.FaceType) (image.Image, error) {
	link2, errReadFile := p.API.ReadFile(currentImagePath)
	if errReadFile != nil {
		return nil, errors.New(errReadFile.Error())
	}

	img, _, errDecodeImage := image.Decode(bytes.NewReader(link2))
	if errDecodeImage != nil {
		return nil, errDecodeImage
	}

	bundlePath, errGetBundlePath := p.API.GetBundlePath()
	if errGetBundlePath != nil {
		return nil, errGetBundlePath
	}

	cascadeFile, errReadCascade := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "facefinder"))
	if errReadCascade != nil {
		return nil, errReadCascade
	}

	return swap.ImgFaceSwap(img, faces, cascadeFile)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Mattermost-User-Id") == "" {
		http.Error(w, "please log in", http.StatusForbidden)
		return
	}

	p.router.ServeHTTP(w, r)
}

// OnActivate activate plugin
func (p *Plugin) OnActivate() error {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return err
	}

	err = facesLib.LoadFacesConfig(bundlePath)

	if err != nil {
		return err
	}

	p.router = mux.NewRouter()
	p.router.HandleFunc("/img/{name}.jpg", serveImg).Methods("GET")

	names := facesLib.Names()

	return p.API.RegisterCommand(&model.Command{
		Trigger:          "faceswap",
		AutoComplete:     true,
		AutoCompleteDesc: strings.Join(names, ", "),
	})
}

// ExecuteCommand face-swap command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/faceswap"))

	faces := facesLib.Faces()
	var memeFaces []*facesLib.FaceType

	if strings.Trim(input, "") == "" {
		memeFaces = faces
	} else {
		selectedFaces := strings.Split(input, ",")

		for index, f := range faces {
			for _, sf := range selectedFaces {
				if f.Name == strings.Trim(sf, " ") {
					memeFaces = append(memeFaces, faces[index])
				}
			}
		}
	}

	if len(currentImagePath) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Before you have to upload an image",
		}, nil
	}

	if len(memeFaces) > 0 {
		id := generateID()

		img, err := p.generateMeme(memeFaces)

		if err != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, err.Error(), http.StatusInternalServerError)
		}

		generatedMemesIds = append(generatedMemesIds, id)

		if len(generatedMemesIds) > maxInMemoryMemes {
			delete(generatedMemes, generatedMemesIds[0])
			generatedMemesIds = generatedMemesIds[1:]
		}

		generatedMemes[id] = img

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
			Text:         "![Face swap](/plugins/faceswap/img/" + id + ".jpg)",
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Valid names: " + strings.Join(facesLib.Names(), ", "),
	}, nil
}

// MessageHasBeenPosted read channel messages
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if len(post.FileIds) > 0 {
		link, err := p.API.GetFileInfo(post.FileIds[0])

		if err != nil {
			fmt.Println(err.Error())
		} else if utils.IsImage(link.Path) {
			currentImagePath = link.Path
		}
	}
}
