package main

import (
	"image"
	"io/ioutil"
	"log"
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
	image image.Image
	name  string
}

var images = []faceType{}

func isImage(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func getImg() {
	files, err := ioutil.ReadDir("./faces")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()

		if !isImage(fileName) {
			// yalm?
			continue
		}

		faceName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		image := loadImage(filepath.Join("faces", fileName))

		face := faceType{
			name:  faceName,
			image: image,
		}

		images = append(images, face)
	}
}
