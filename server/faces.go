package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func loadImage(filePath string) image.Image {
	img, err := imaging.Open(filePath)
	if err != nil {
		panic(err)
	}
	return img
}

type faceType struct {
	image       image.Image
	name        string
	width       float32
	height      float32
	paddingLeft float32
	paddingTop  float32
}

// FaceSwapConfig json file config
type FaceSwapConfig struct {
	Faces []struct {
		Name   string `json:"name"`
		Images []struct {
			Path        string  `json:"path"`
			Width       float32 `json:"width,omitempty"`
			Height      float32 `json:"height,omitempty"`
			PaddingLeft float32 `json:"paddingLeft,omitempty"`
			PaddingTop  float32 `json:"paddingTop,omitempty"`
		} `json:"images"`
	} `json:"faces"`
	Width       float32 `json:"width"`
	Height      float32 `json:"height"`
	PaddingLeft float32 `json:"paddingLeft"`
	PaddingTop  float32 `json:"paddingTop"`
}

var images = []faceType{}

func isImage(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func loadImgConfigFile(filePath string) FaceSwapConfig {
	var faceSwapConfig FaceSwapConfig

	jsonFile, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &faceSwapConfig)

	fmt.Println(faceSwapConfig)

	for _, data := range faceSwapConfig.Faces {
		fmt.Println("name: " + data.Name)
	}

	defer jsonFile.Close()

	return faceSwapConfig
}

func loadImages(bundlePath string) {
	filePath := filepath.Join(bundlePath, "assets", "faces.json")
	var faceSwapConfig = loadImgConfigFile(filePath)

	for _, data := range faceSwapConfig.Faces {
		fmt.Println("name: " + data.Name)

		for _, image := range data.Images {
			imageFile := loadImage(filepath.Join(bundlePath, "assets", image.Path))
			width := image.Width
			height := image.Height
			paddingLeft := image.PaddingLeft
			paddingTop := image.PaddingTop

			if width == 0 {
				width = faceSwapConfig.Width
			}
			if height == 0 {
				height = faceSwapConfig.Height
			}
			if paddingLeft == 0 {
				paddingLeft = faceSwapConfig.PaddingLeft
			}
			if paddingTop == 0 {
				paddingTop = faceSwapConfig.PaddingTop
			}

			newFace := faceType{
				image:       imageFile,
				name:        data.Name,
				width:       width,
				height:      height,
				paddingLeft: paddingLeft,
				paddingTop:  paddingTop,
			}

			images = append(images, newFace)
		}
	}

	fmt.Println(images)
}

// Faces get faces
func Faces() []faceType {
	return images
}

func main() {
	loadImages("../")
}

// func getImg() {
// 	files, err := ioutil.ReadDir("./faces")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		fileName := file.Name()

// 		if !isImage(fileName) {
// 			// yalm?
// 			continue
// 		}

// 		faceName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
// 		image := loadImage(filepath.Join("faces", fileName))

// 		face := faceType{
// 			name:  faceName,
// 			image: [image],
// 		}

// 		images = append(images, face)
// 	}
// }
