package utils

import (
	"image"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// AppendIfMissing append if is not in the array
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// LoadImage return the path image
func LoadImage(filePath string) image.Image {
	img, err := imaging.Open(filePath)
	if err != nil {
		panic(err)
	}
	return img
}

// IsImage check if the file is an image
func IsImage(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
