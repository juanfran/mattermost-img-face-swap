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

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/segmentio/ksuid"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
	botUserID     string
	lastImagePath map[string]string
}

func generateID() string {
	return ksuid.New().String()
}

func (p *Plugin) generateMeme(currentImagePath string, faces []*facesLib.FaceType) (image.Image, error) {
	link, errReadFile := p.API.ReadFile(currentImagePath)
	if errReadFile != nil {
		return nil, errors.New(errReadFile.Error())
	}

	img, _, errDecodeImage := image.Decode(bytes.NewReader(link))
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

// OnActivate activate plugin
func (p *Plugin) OnActivate() error {
	p.lastImagePath = map[string]string{}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return err
	}

	err = facesLib.LoadFacesConfig(bundlePath)

	if err != nil {
		return err
	}

	bot := &model.Bot{
		Username:    "faceswap",
		DisplayName: "Faceswap",
	}
	botUserID, ensureBotErr := p.Helpers.EnsureBot(bot)
	if ensureBotErr != nil {
		return ensureBotErr
	}
	p.botUserID = botUserID

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

	if _, ok := p.lastImagePath[args.ChannelId]; !ok {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Before you have to upload an image",
		}, nil
	}

	if len(memeFaces) > 0 {
		id := generateID()
		img, err := p.generateMeme(p.lastImagePath[args.ChannelId], memeFaces)
		if err != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, err.Error(), http.StatusInternalServerError)
		}

		buf := new(bytes.Buffer)
		jpeg.Encode(buf, img, nil)
		imgBytes := buf.Bytes()

		imgInfo, errUpload := p.API.UploadFile(imgBytes, args.ChannelId, id+".jpg")

		if errUpload != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, errUpload.Error(), http.StatusInternalServerError)
		}

		var fileIDs []string
		fileIDs = append(fileIDs, imgInfo.Id)

		post := &model.Post{
			UserId:    p.botUserID,
			ChannelId: args.ChannelId,
			RootId:    args.RootId,
			FileIds:   fileIDs,
		}

		_, createPostError := p.API.CreatePost(post)

		if createPostError != nil {
			return nil, model.NewAppError("ExecuteCommand", "error faceswap", nil, createPostError.Error(), http.StatusInternalServerError)
		}

		return &model.CommandResponse{}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Valid names: " + strings.Join(facesLib.Names(), ", "),
	}, nil
}

// MessageHasBeenPosted read channel messages
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if len(post.FileIds) > 0 && post.UserId != p.botUserID {
		link, err := p.API.GetFileInfo(post.FileIds[0])

		if err != nil {
			fmt.Println(err.Error())
		} else if link.IsImage() {
			p.lastImagePath[post.ChannelId] = link.Path
		}
	}
}
