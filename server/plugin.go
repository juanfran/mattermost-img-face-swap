package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"image/jpeg"
	_ "image/jpeg"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/segmentio/ksuid"
)

// MaxInMemoryMemes how many meme will stores in memory
var MaxInMemoryMemes = 10

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	// configuration *configuration

	router *mux.Router
}

// GeneratedMemes store images
var GeneratedMemes map[string]image.Image = make(map[string]image.Image)

var generatedMemesIds []string
var currentImagePath string

// generateID, generate random id with ksuid
func generateID() string {
	return ksuid.New().String()
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Mattermost-User-Id") == "" {
		http.Error(w, "please log in", http.StatusForbidden)
		return
	}

	p.router.ServeHTTP(w, r)
}

// serve Img
func serveImg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=604800")

	vars := mux.Vars(r)
	memeID := vars["name"]

	if err := jpeg.Encode(w, GeneratedMemes[memeID], &jpeg.Options{
		Quality: 90,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// OnActivate activate plugin
func (p *Plugin) OnActivate() error {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		panic(err)
	}
	loadImages(bundlePath)

	p.router = mux.NewRouter()
	p.router.HandleFunc("/img/{name}.jpg", serveImg).Methods("GET")

	return p.API.RegisterCommand(&model.Command{
		Trigger:          "faceswap",
		AutoComplete:     true,
		AutoCompleteDesc: "todo2.",
	})
}

// ExecuteCommand face-swap command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/faceswap"))

	faces := Faces()
	var memeFaces []FaceType

	if strings.Trim(input, "") == "" {
		memeFaces = faces
	} else {
		selectedFaces := strings.Split(input, ",")

		for index, f := range faces {
			for _, sf := range selectedFaces {
				if f.name == strings.Trim(sf, " ") {
					memeFaces = append(memeFaces, faces[index])
				}
			}
		}
	}

	if len(memeFaces) > 0 {
		id := generateID()

		img, err := p.generateMeme(memeFaces)

		if err != nil {
			fmt.Printf("todo error 2")
		}

		generatedMemesIds = append(generatedMemesIds, id)

		if len(generatedMemesIds) > MaxInMemoryMemes {
			delete(GeneratedMemes, generatedMemesIds[0])
			generatedMemesIds = generatedMemesIds[1:]
		}

		GeneratedMemes[id] = img

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
			Text:         "![Face swap](/plugins/faceswap/img/" + id + ".jpg)",
		}, nil
	}

	var names []string

	for _, f := range faces {
		names = append(names, f.name)
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Valid names:" + strings.Join(names, ", "),
	}, nil
}

// MessageHasBeenPosted read channel messages
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if len(post.FileIds) > 0 {
		link, err := p.API.GetFileInfo(post.FileIds[0])
		if err != nil {
			fmt.Printf("todo error 1")
		}

		currentImagePath = link.Path
	}

}

// GenerateMeme generate meme with the last image file
func (p *Plugin) generateMeme(faces []FaceType) (image.Image, error) {
	link2, err2 := p.API.ReadFile(currentImagePath)
	if err2 != nil {
		fmt.Printf("todo error 2")
	}

	img, _, err3 := image.Decode(bytes.NewReader(link2))
	if err3 != nil {
		fmt.Printf("todo error 3")
	}

	bundlePath, errBundlePath := p.API.GetBundlePath()
	if errBundlePath != nil {
		fmt.Printf("todo error 4")
	}

	cascadeFile, errCascadeFile := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "facefinder"))

	if errCascadeFile != nil {
		fmt.Printf("Error reading the cascade file: %v", errCascadeFile)
	}

	return faceswap(img, faces, cascadeFile)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
