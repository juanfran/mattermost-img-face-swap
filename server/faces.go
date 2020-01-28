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

// FaceType loaded from config
type FaceType struct {
	image       string
	name        string
	width       int
	paddingLeft int
	paddingTop  int
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

var images = []FaceType{}

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
	fmt.Printf("filePath------>: %v \n", filePath)
	var faceSwapConfig = loadImgConfigFile(filePath)

	for _, data := range faceSwapConfig.Faces {
		fmt.Println("name: " + data.Name)

		for _, image := range data.Images {
			imageFile := filepath.Join(bundlePath, "assets", image.Path)

			// src, err := pigo.GetImage(imageFile)
			// fmt.Printf("444444444444444444444444444444444444444444444444444")
			// if err != nil {
			// 	fmt.Println("NOOOOOOOOOOOOO")
			// } else {
			// 	fmt.Println("SIIIIIIIIIIIIIII")
			// 	fmt.Println(src)

			// }

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

			newFace := FaceType{
				image:       imageFile,
				name:        data.Name,
				width:       width,
				paddingLeft: paddingLeft,
				paddingTop:  paddingTop,
			}

			images = append(images, newFace)
		}
	}

	fmt.Println(images)
	fmt.Printf("Len: %v", len(images))
	fmt.Printf("Name: %v", images[0].name)

	// test := filepath.Join(bundlePath, "assets", "meme3.jpg")
	// faceswap(test, images[0])
}

// Faces get faces
func Faces() []FaceType {
	return images
}

/* func main() {
	loadImages("../")
} */

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
// 			image: image,
// 		}

// 		images = append(images, face)
// 	}
// }
