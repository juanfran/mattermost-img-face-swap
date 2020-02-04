package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

func faceswap(image1 image.Image, memeFace FaceType, cascadeFile []byte) (image.Image, error) {
	image2, err := imaging.Open(memeFace.image)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	src := pigo.ImgToNRGBA(image1)

	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	pixels := pigo.RgbToGrayscale(src)

	fmt.Printf("len: %v\n", len(pixels))

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

	fmt.Printf("cParams: %v\n", cParams.MinSize)
	fmt.Printf("cascadeFile: %v\n", len(cascadeFile))

	pigo := pigo.NewPigo()
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	angle := 0.0
	dets := classifier.RunCascade(cParams, angle)
	dets = classifier.ClusterDetections(dets, 0.2)
	var final image.Image

	for _, face := range dets {
		x := float64(face.Col - face.Scale/2)
		y := float64(face.Row - face.Scale/2)
		width := face.Scale * (100 / memeFace.width)
		fmt.Printf("width: %v\n", width)
		resizedImage := imaging.Resize(image2, width, 0, imaging.Lanczos)

		bounds := resizedImage.Bounds()
		fmt.Println("Width Real:", bounds.Max.X, "Height:", bounds.Max.Y)

		dst := imaging.New(image1.Bounds().Dx(), image1.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
		dst = imaging.Paste(dst, image1, image.Pt(0, 0))
		position := image.Pt(
			int(x)+memeFace.paddingLeft,
			(int(y)*80/100)+memeFace.paddingTop)

		final = imaging.Overlay(dst, resizedImage, position, 1)
	}

	return final, nil
}
