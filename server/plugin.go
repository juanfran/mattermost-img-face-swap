package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

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

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// serve Img
func serveImg(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	return
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
		Text:         "![Face swap](/plugins/face-swap/img/test.jpg)",
	}, nil
}

// MessageHasBeenPosted read channel messages
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	fmt.Printf("---------------------------------------------")
	fmt.Printf("Post: %v\n", post.Message)
	fmt.Printf("Post FileIds: %v\n", post.FileIds)

	if len(post.FileIds) > 0 {
		fmt.Printf("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeee: \n")
		link, err := p.API.GetFile(post.FileIds[0])
		fmt.Printf("Liiiiiiiiiiiiiink: %v\n", post.FileIds[0])
		fmt.Printf("Liiiiiiiiiiiiiink: %v\n", err)

		if err == nil {
			fmt.Printf("Post Link: %v\n", link)
		}
	}

}

// See https://developers.mattermost.com/extend/plugins/server/reference/
