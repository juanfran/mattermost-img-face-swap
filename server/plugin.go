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
)

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

var currentImage image.Image = nil

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Mattermost-User-Id") == "" {
		http.Error(w, "please log in", http.StatusForbidden)
		return
	}

	p.router.ServeHTTP(w, r)
}

// serve Img
func serveImg(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SERVE IMG!! \n")
	/* 	http.NotFound(w, r)
	   	return */

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=604800")

	fmt.Printf("serveImg serveImg serveImg serveImg serveImg serveImg serveImg \n")

	if currentImage != nil {
		fmt.Printf("currentImage != nil")
	} else {
		fmt.Printf("currentImage == nil")
	}
	// http.NotFound(w, r)
	// return

	if err := jpeg.Encode(w, currentImage, &jpeg.Options{
		Quality: 90,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// OnActivate activate plugin
func (p *Plugin) OnActivate() error {
	fmt.Printf("activate activate activate \n")

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
		AutoCompleteDesc: "todo.",
	})
}

// ExecuteCommand face-swap command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/faceswap"))
	fmt.Printf("Input2: %v\n", input)

	names := []string{"person1", "person2"}

	if input == "" {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Names:" + strings.Join(names, ", "),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         "![Face swap](/plugins/faceswap/img/test.jpg)",
	}, nil
}

// MessageHasBeenPosted read channel messages
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	fmt.Printf("---------------------------------------------")
	fmt.Printf("Post: %v\n", post.Message)
	fmt.Printf("Post FileIds: %v\n", post.FileIds)

	// /mm/mattermost/plugins/faceswap/assets/faces.json
	// /mm/mattermost-data/mm/mattermost/plugins/faceswap/assets/facefinder

	if len(post.FileIds) > 0 {
		fmt.Printf("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeee: \n")
		/* 		link, err := p.API.GetFile(post.FileIds[0])
		   		fmt.Printf("Liiiiiiiiiiiiiink: %v\n", post.FileIds[0])
		   		fmt.Printf("Liiiiiiiiiiiiiink: %v\n", err) */

		link, err := p.API.GetFileInfo(post.FileIds[0])
		if err == nil {
			fmt.Printf("todo error 1")
		}

		link2, err2 := p.API.ReadFile(link.Path)
		if err2 == nil {
			fmt.Printf("todo error 2")
		}

		img, _, err3 := image.Decode(bytes.NewReader(link2))
		if err3 == nil {
			fmt.Printf("todo error 3")
		}

		if err == nil {
			fmt.Printf("Post Path: %v\n", link.Path)
		}

		// fmt.Println(Faces()[0].image)
		bundlePath, errBundlePath := p.API.GetBundlePath()
		if errBundlePath != nil {
			panic(err)
		}

		// cascadeFile, err := p.API.ReadFile(filepath.Join(bundlePath, "assets", "facefinder"))

		cascadeFile, errCascadeFile := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "facefinder"))

		fmt.Printf("cascadeFile: %v\n", len(cascadeFile))

		if errCascadeFile != nil {
			fmt.Printf("Error reading the cascade file: %v", err)
		}

		resultImage, imgError := faceswap(img, Faces()[0], cascadeFile)

		currentImage = resultImage

		if imgError != nil {
			fmt.Printf("Error faceswaping: %v\n", link.Path)
		}
	}

}

// import (
// 	"bytes"
// 	"image"
// 	_ "image/jpeg"
// 	_ "image/png"
// )

// func mustLoadImage(assetName string) image.Image {
// 	img, _, err := image.Decode(bytes.NewReader(MustAsset(assetName)))
// 	if err != nil {
// 		panic(err)
// 	}
// 	return img
// }

// See https://developers.mattermost.com/extend/plugins/server/reference/
