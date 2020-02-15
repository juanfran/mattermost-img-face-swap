package faces

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juanfran/mattermost-img-face-swap/server/utils"
)

// FaceType loaded from config
type FaceType struct {
	Image       string
	Name        string
	Width       int
	PaddingLeft int
	PaddingTop  int
}

// FaceSwapConfig json file config
type FaceSwapConfig struct {
	Faces []struct {
		Name   string `json:"name"`
		Images []struct {
			Path        string `json:"path"`
			Width       int    `json:"width,omitempty"`
			PaddingLeft int    `json:"paddingLeft,omitempty"`
			PaddingTop  int    `json:"paddingTop,omitempty"`
		} `json:"images"`
	} `json:"faces"`
	Width       int `json:"width"`
	PaddingLeft int `json:"paddingLeft"`
	PaddingTop  int `json:"paddingTop"`
}

var images = []*FaceType{}

func loadImgConfigFile(filePath string) (*FaceSwapConfig, error) {
	var faceSwapConfig FaceSwapConfig

	jsonFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &faceSwapConfig)

	defer jsonFile.Close()

	return &faceSwapConfig, nil
}

// LoadFacesConfig load faces config
func LoadFacesConfig(bundlePath string) error {
	filePath := filepath.Join(bundlePath, "assets", "faces.json")
	faceSwapConfig, err := loadImgConfigFile(filePath)

	if err != nil {
		return err
	}

	for _, data := range faceSwapConfig.Faces {
		for _, image := range data.Images {
			imageFile := filepath.Join(bundlePath, "assets", image.Path)

			width := image.Width
			paddingLeft := image.PaddingLeft
			paddingTop := image.PaddingTop

			if width == 0 {
				width = faceSwapConfig.Width
			}
			if paddingLeft == 0 {
				paddingLeft = faceSwapConfig.PaddingLeft
			}
			if paddingTop == 0 {
				paddingTop = faceSwapConfig.PaddingTop
			}

			newFace := &FaceType{
				Image:       imageFile,
				Name:        data.Name,
				Width:       width,
				PaddingLeft: paddingLeft,
				PaddingTop:  paddingTop,
			}

			images = append(images, newFace)
		}
	}

	return nil
}

// Faces get faces
func Faces() []*FaceType {
	return images
}

// Names return all face names
func Names() []string {
	var names []string

	for _, f := range Faces() {
		names = utils.AppendIfMissing(names, f.Name)
	}

	return names
}
