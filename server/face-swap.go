package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
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
		MinSize:     50,
		MaxSize:     640,
		ShiftFactor: 0.2,
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

	final = image1

	for index, face := range dets {
		if index < len(memeFaces) {
			memeFace := memeFaces[index]
			fmt.Println(memeFace.name)

			image2, err := imaging.Open(memeFace.image)

			if err != nil {
				log.Fatalf("failed to open image: %v", err)
			}

			col := float64(face.Col)
			row := float64(face.Row)
			scale := float64(face.Scale)

			x := int(col - scale/2)
			y := int(row - scale/2)

			width := int(math.Round(scale * (float64(memeFace.width) / 100.00)))

			fmt.Println("width width width width width")
			fmt.Println(scale)
			fmt.Println(memeFace.width)
			fmt.Println(width)

			resizedImage := imaging.Resize(image2, width, 0, imaging.Lanczos)
			resizedImageBounds := resizedImage.Bounds()

			dst := imaging.New(final.Bounds().Dx(), final.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
			dst = imaging.Paste(dst, final, image.Pt(0, 0))
			position := image.Pt(
				x+int(float64(resizedImageBounds.Size().X)*float64(memeFace.paddingLeft)/100.00),
				y+int(float64(resizedImageBounds.Size().Y)*float64(memeFace.paddingTop)/100.00),
			)

			final = imaging.Overlay(dst, resizedImage, position, 1)
		}
	}

	return final, nil
}
