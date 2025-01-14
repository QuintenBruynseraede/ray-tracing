package internal

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
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

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}
