package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

// RgbToGrayscale test
func RgbToGrayscale(src image.Image) []uint8 {
	cols, rows := src.Bounds().Dx(), src.Bounds().Dy()
	gray := make([]uint8, rows*cols)

	// todo debug outofrange
	fmt.Printf("cols: %v\n", cols)
	fmt.Printf("rows: %v\n", rows)
	fmt.Printf("05")

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			// src.At(x, y).RGBA()
			// r, g, b, _ := src.At(x, y).RGBA()
			// gray[y*cols+x] = uint8(
			// 	(0.299*float64(r) +
			// 		0.587*float64(g) +
			// 		0.114*float64(b)) / 256,
			// )
		}
	}
	return gray
}

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
		MinSize:     20,
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
		width := face.Scale
		resizedImage := imaging.Resize(image2, width, 0, imaging.Lanczos)

		dst := imaging.New(image1.Bounds().Dx(), image1.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
		dst = imaging.Paste(dst, image1, image.Pt(0, 0))
		final = imaging.Overlay(dst, resizedImage, image.Pt(int(x)+memeFace.paddingLeft, (int(y)*80/100)+memeFace.paddingTop), 1)
		// err = imaging.Save(dst, "result.jpg")
		// if err != nil {
		// 	log.Fatalf("failed to save image: %v", err)
		// }
	}

	return final, nil
}

// image1, err := imaging.Open(memePath)
// if err != nil {
// 	log.Fatalf("failed to open image: %v", err)
// }

// image2, err := imaging.Open(face.image)
// if err != nil {
// 	log.Fatalf("failed to open image: %v", err)
// }

// fmt.Println(memePath)
// fmt.Println(face.image)

// // load classifier to recognize faces
// classifier := gocv.NewCascadeClassifier()
// defer classifier.Close()

// xmlFile := "./haarcascade_frontalface_alt.xml"
// if !classifier.Load(xmlFile) {
// 	fmt.Printf("Error reading cascade file: %v\n", xmlFile)
// 	return
// }

// // detect faces
// rects := classifier.DetectMultiScale(bigImage)

// for _, r := range rects {
// 	// fmt.Printf("Max.X %v\n", r.Max.X)
// 	// fmt.Printf("Max.Y %v\n", r.Max.Y)
// 	// fmt.Printf("Min.X %v\n", r.Min.X)
// 	// fmt.Printf("Min.Y %v\n", r.Min.Y)
// 	width := r.Max.X - r.Min.X
// 	// fmt.Printf("width: %v\n", width)
// 	// fmt.Printf("height: %v\n", height)
// 	newWidth := width + int(face.width*(width/100)/2)
// 	resizedImage := imaging.Resize(image2, newWidth, 0, imaging.Lanczos)

// 	dst := imaging.New(image1.Bounds().Dx(), image1.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
// 	dst = imaging.Paste(dst, image1, image.Pt(0, 0))
// 	dst = imaging.Overlay(dst, resizedImage, image.Pt(r.Min.X+face.paddingLeft, (r.Min.Y-r.Min.Y*80/100)+face.paddingTop), 1)

// 	// return dst
// 	// image.NRGBA
// 	err = imaging.Save(dst, "result.jpg")
// 	if err != nil {
// 		log.Fatalf("failed to save image: %v", err)
// 	}
// }
