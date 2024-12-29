package internal

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"time"
)

// SaveScreenshot saves an image to to a folder out/.
// The file name is "out-$TIMESTAMP.png".
func SaveScreenshot(image *image.RGBA) {
	fileName := fmt.Sprintf("out-%s.jpg", time.Now().Format("20060102-150405"))
	err := os.Mkdir("./out", os.ModePerm)
	if !errors.Is(err, os.ErrExist) {
		panic(err)
	}

	f, err := os.Create("./out/" + fileName)
	if err != nil {
		panic(err)
	}

	if err := jpeg.Encode(f, image, &jpeg.Options{Quality: 100}); err != nil {
		f.Close()
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
	log.Println("Done")
}

func Gradient(primary color.RGBA, secondary color.RGBA, t float64) color.RGBA {
	if t < 0.0 || t > 1.0 {
		log.Fatal("t must be between 0.0 and 1.0")
	}
	return color.RGBA{
		R: uint8(t*float64(primary.R) + (1-t)*float64(secondary.R)),
		G: uint8(t*float64(primary.G) + (1-t)*float64(secondary.G)),
		B: uint8(t*float64(primary.B) + (1-t)*float64(secondary.B)),
		A: uint8(t*float64(primary.A) + (1-t)*float64(secondary.A)),
	}
}

func AddColor(a color.RGBA, b color.RGBA) color.RGBA {
	return color.RGBA{
		R: a.R + b.R,
		G: a.G + b.G,
		B: a.B + b.B,
		A: a.A,
	}
}

func MultiplyColor(a color.RGBA, x float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(a.R) * x),
		G: uint8(float64(a.G) * x),
		B: uint8(float64(a.B) * x),
		A: a.A,
	}
}
