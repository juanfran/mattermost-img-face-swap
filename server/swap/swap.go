package swap

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	facesLib "github.com/juanfran/mattermost-img-face-swap/server/faces"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

func shuffleFaces(a []*facesLib.FaceType) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

// ImgFaceSwap swap image faces by newFaces
func ImgFaceSwap(sourceImg image.Image, newFaces []*facesLib.FaceType, cascadeFile []byte) (image.Image, error) {
	src := pigo.ImgToNRGBA(sourceImg)

	shuffleFaces(newFaces)

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
	final := sourceImg

	for index, face := range dets {
		if index < len(newFaces) {
			newFace := newFaces[index]

			newFaceImage, err := imaging.Open(newFace.Image)

			if err != nil {
				log.Fatalf("failed to open image: %v", err)
			}

			col := float64(face.Col)
			row := float64(face.Row)
			scale := float64(face.Scale)

			x := int(col - scale/2)
			y := int(row - scale/2)

			width := int(math.Round(scale * (float64(newFace.Width) / 100.00)))

			resizedImage := imaging.Resize(newFaceImage, width, 0, imaging.Lanczos)
			resizedImageBounds := resizedImage.Bounds()

			dst := imaging.New(final.Bounds().Dx(), final.Bounds().Dy(), color.NRGBA{0, 0, 0, 0})
			dst = imaging.Paste(dst, final, image.Pt(0, 0))

			position := image.Pt(
				x+int(float64(resizedImageBounds.Size().X)*float64(newFace.PaddingLeft)/100.00),
				y+int(float64(resizedImageBounds.Size().Y)*float64(newFace.PaddingTop)/100.00),
			)

			final = imaging.Overlay(dst, resizedImage, position, 1)
		}
	}

	return final, nil
}
