package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

// Shuffle face type
func Shuffle(a []FaceType) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func faceswap(image1 image.Image, memeFaces []FaceType, cascadeFile []byte) (image.Image, error) {
	src := pigo.ImgToNRGBA(image1)

	Shuffle(memeFaces)

	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	pixels := pigo.RgbToGrayscale(src)

	cParams := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     640,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	pigo := pigo.NewPigo()
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	angle := 0.0
	dets := classifier.RunCascade(cParams, angle)
	dets = classifier.ClusterDetections(dets, 0.2)
	var final image.Image

	fmt.Println("--------------")
	fmt.Println(len(memeFaces))
	fmt.Println(len(dets))

	final = image1

	for index, face := range dets {
		if index < len(memeFaces) {
			memeFace := memeFaces[index]
			fmt.Println(memeFace.name)

			image2, err := imaging.Open(memeFace.image)

			if err != nil {
				log.Fatalf("failed to open image: %v", err)
			}

			fmt.Println("Image loaded")

			x := float64(face.Col - face.Scale/2)
			y := float64(face.Row - face.Scale/2)
			width := face.Scale * (100 / memeFace.width)

			resizedImage := imaging.Resize(image2, width, 0, imaging.Lanczos)

			dst := imaging.New(final.Bounds().Dx(), final.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
			dst = imaging.Paste(dst, final, image.Pt(0, 0))
			position := image.Pt(
				int(x)+memeFace.paddingLeft,
				int(y)+memeFace.paddingTop)

			final = imaging.Overlay(dst, resizedImage, position, 1)
		}
	}

	return final, nil
}
