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
	images      []image.Image
	name        string
	width       string
	height      string
	paddingLeft string
	paddintTop  string
}

type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
	Width       int    `json:"width, omitempty"`
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

func loadImgConfigFile(filePath string) {
	// var faceSwapConfig map[string]interface{}
	var faceSwapConfig FaceSwapConfig

	jsonFile, err := os.Open(filePath)

	fmt.Println(filePath)

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
}

func loadImages(bundlePath string) {
	dir := filepath.Join(bundlePath, "assets", "faces")
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()

		ext := filepath.Ext(fileName)
		if ext != ".json" {
			continue
		}

		loadImgConfigFile(filepath.Join(dir, file.Name()))
	}
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
