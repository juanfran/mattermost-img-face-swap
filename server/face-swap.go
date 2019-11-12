package main

import "C"
import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	"gocv.io/x/gocv"
)

func faceswap(memePath string, face faceType) {
	bigImage := gocv.IMRead(memePath, gocv.IMReadColor)

	image1, err := imaging.Open(memePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	image2, err := imaging.Open(face.image)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	fmt.Println(memePath)
	fmt.Println(face.image)

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	// xmlFile := os.Args[1]
	xmlFile := "/usr/share/opencv4/haarcascades/haarcascade_frontalface_alt.xml"
	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	// detect faces
	rects := classifier.DetectMultiScale(bigImage)
	fmt.Printf("found %d faces\n", len(rects))

	for _, r := range rects {
		fmt.Printf("Max.X %v\n", r.Max.X)
		fmt.Printf("Max.Y %v\n", r.Max.Y)
		fmt.Printf("Min.X %v\n", r.Min.X)
		fmt.Printf("Min.Y %v\n", r.Min.Y)
		width := r.Max.X - r.Min.X
		height := r.Max.Y - r.Min.Y

		fmt.Printf("width: %v\n", width)
		fmt.Printf("height: %v\n", height)
		newWidth := width + int(30*(width/100)/2)

		resizedImage := imaging.Resize(image2, newWidth, 0, imaging.Lanczos)

		dst := imaging.New(image1.Bounds().Dx(), image1.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
		dst = imaging.Paste(dst, image1, image.Pt(0, 0))
		dst = imaging.Overlay(dst, resizedImage, image.Pt(r.Min.X, r.Min.Y-r.Min.Y*80/100), 1)

		// return dst
		// image.NRGBA
		err = imaging.Save(dst, "result.jpg")
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
	}
}
