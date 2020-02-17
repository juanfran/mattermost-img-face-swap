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
	fmt.Println("aaaaaa")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=604800")

	vars := mux.Vars(r)
	memeID := vars["name"]
	fmt.Println("bbbb")
	if err := jpeg.Encode(w, generatedMemes[memeID], &jpeg.Options{
		Quality: 90,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) generateMeme(faces []*facesLib.FaceType) (image.Image, error) {
	fmt.Println("6666")
	link2, errReadFile := p.API.ReadFile(currentImagePath)
	if errReadFile != nil {
		return nil, errors.New(errReadFile.Error())
	}
	fmt.Println("777")
	img, _, errDecodeImage := image.Decode(bytes.NewReader(link2))
	if errDecodeImage != nil {
		return nil, errDecodeImage
	}
	fmt.Println("888")
	bundlePath, errGetBundlePath := p.API.GetBundlePath()
	if errGetBundlePath != nil {
		return nil, errGetBundlePath
	}
	fmt.Println("999")
	cascadeFile, errReadCascade := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "facefinder"))
	if errReadCascade != nil {
		return nil, errReadCascade
	}
	fmt.Println("00000")
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
		fmt.Println("11111")
		id := generateID()
		fmt.Println("222")
		img, err := p.generateMeme(memeFaces)
		fmt.Println("333")
		if err != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, err.Error(), http.StatusInternalServerError)
		}

		generatedMemesIds = append(generatedMemesIds, id)
		fmt.Println("4444")
		if len(generatedMemesIds) > maxInMemoryMemes {
			delete(generatedMemes, generatedMemesIds[0])
			generatedMemesIds = generatedMemesIds[1:]
		}

		generatedMemes[id] = img

		buf := new(bytes.Buffer)
		jpeg.Encode(buf, img, nil)
		imgBytes := buf.Bytes()

		fmt.Println("555")
		fmt.Printf("image path: %v", "/plugins/faceswap/img/"+id+".jpg")

		imgInfo, errUpload := p.API.UploadFile(imgBytes, args.ChannelId, id+".jpg")

		if errUpload != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, errUpload.Error(), http.StatusInternalServerError)
		}

		var fileIDs []string
		fileIDs = append(fileIDs, imgInfo.Id)

		post := &model.Post{
			UserId:    args.UserId,
			ChannelId: args.ChannelId,
			RootId:    args.RootId,
			FileIds:   fileIDs,
		}

		_, createPostError := p.API.CreatePost(post)

		if createPostError != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, createPostError.Error(), http.StatusInternalServerError)
		}

		// p.API.SendEphemeralPost(args.UserId, post)

		return &model.CommandResponse{}, nil
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
			fmt.Println("NEW IMAGE")
			fmt.Println(link.Path)
			currentImagePath = link.Path
		}
	}
}

// ,"msg":"NEW IMAGE","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"20200217/teams/noteam/channels/eeyrbxyqu3dduy7ukm99uqepyo/users/9zbb9jj73igzdxz5cy4je7xa9y/xphospotqidkbp34z9fmpeon1y/meme2.jpg","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"11111","plugin_id":"faceswap","source":"plugin_stdout"}
// "msg":"222","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"6666","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"777","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"888","plugin_id":"faceswap","source":"plugin_stdout"}
// "msg":"999","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"00000","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"333","plugin_id":"faceswap","source":"plugin_stdout"}
// ,"msg":"4444","plugin_id":"faceswap","source":"plugin_stdout"}
// "msg":"555","plugin_id":"faceswap","source":"plugin_stdout"}
