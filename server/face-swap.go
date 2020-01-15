package main

import (
	"image"
	"image/color"
	"log"

	"io/ioutil"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

func faceswap(memePath string, memeFace FaceType) (image.Image, error) {
	image1, err := imaging.Open(memePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	image2, err := imaging.Open(memeFace.image)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	cascadeFile, err := ioutil.ReadFile("./facefinder")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	src, err := pigo.GetImage(memePath)
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

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
