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

// FaceSwapConfig json file config
type FaceSwapConfig struct {
	faces []struct {
		name   string
		images []struct {
			path        string
			width       string
			height      string
			paddingLeft string
			paddintTop  string
		}
	}
	width       string
	height      string
	paddingLeft string
	paddintTop  string
}

var images = []faceType{}

func isImage(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func loadImgConfigFile(filePath string) {
	var faceSwapConfig FaceSwapConfig

	jsonFile, err := os.Open(filePath)

	fmt.Println(filePath)

	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &faceSwapConfig)

	fmt.Println(faceSwapConfig)

	for i := 0; i < len(faceSwapConfig.faces); i++ {
		fmt.Println("name: " + faceSwapConfig.faces[i].name)
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
