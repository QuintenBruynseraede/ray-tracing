package internal

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"
)

const (
	VIEWPORT_HEIGHT     = 2.0
	CAMERA_FOCAL_LENGTH = 1.0
	SAMPLES_PER_PIXEL   = 10
	SAMPLE_SCALE        = 1.0 / SAMPLES_PER_PIXEL
)

type Camera struct {
	Center      *Vec3
	Viewport    Viewport
	imageWidth  int
	imageHeight int
}

func NewCamera(center *Vec3, imageWidth, imageHeight int) Camera {
	viewPortWidth := VIEWPORT_HEIGHT * (float64(imageWidth) / float64(imageHeight))
	viewportU := NewVec3(viewPortWidth, 0, 0)
	viewportV := NewVec3(0, -VIEWPORT_HEIGHT, 0)

	// Calculate horizontal and vertical delta vectors from pixel to pixel
	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	// Calculate the location of the upper left pixel
	viewportUpperLeft := (center.
		Sub(NewVec3(0, 0, CAMERA_FOCAL_LENGTH)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2)))

	// Calculate the location of the lower right pixel
	pixel00Location := viewportUpperLeft.
		Add(pixelDeltaU.Mul(0.5)).
		Add(pixelDeltaV.Mul(0.5))

	viewport := Viewport{
		U:               viewportU,
		V:               viewportV,
		PixelDeltaU:     pixelDeltaU,
		PixelDeltaV:     pixelDeltaV,
		UpperLeft:       viewportUpperLeft,
		Pixel00Location: pixel00Location,
	}
	return Camera{
		Center:      center,
		Viewport:    viewport,
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
	}
}

func (c *Camera) Render(image *image.RGBA, world *HittableList) *image.RGBA {
	for y := 0; y < c.imageHeight; y++ {
		for x := 0; x < c.imageWidth; x++ {
			// Use normal ints for intermediate color samples to allow values > 255
			r, g, b := 0, 0, 0
			for sample := 0; sample < SAMPLES_PER_PIXEL; sample++ {
				sampleColor := c.rayColor(c.getRay(x, y), world)
				r += int(sampleColor.R)
				g += int(sampleColor.G)
				b += int(sampleColor.B)
			}
			pixelColor := color.RGBA{
				R: uint8(float64(r) * SAMPLE_SCALE),
				G: uint8(float64(g) * SAMPLE_SCALE),
				B: uint8(float64(b) * SAMPLE_SCALE),
			}
			image.Set(x, y, pixelColor)
		}
	}
	return image
}

func (c *Camera) getRay(x, y int) *Ray {
	offsetX, offsetY := sampleSquare()
	pixelSample := c.Viewport.Pixel00Location.
		Add(c.Viewport.PixelDeltaU.Mul(float64(x) + offsetX)).
		Add(c.Viewport.PixelDeltaV.Mul(float64(y) + offsetY))
	return NewRay(c.Center, pixelSample.Sub(c.Center))
}

func sampleSquare() (float64, float64) {
	return rand.Float64() - 0.5, rand.Float64() - 0.5
}

func (c *Camera) rayColor(ray *Ray, world *HittableList) color.RGBA {
	intensity := Interval{Min: 0, Max: 0.999}
	hit, hitRecord := world.Hit(ray, &Interval{0, math.MaxFloat64})

	if hit {
		N := hitRecord.N.Normalize()
		return color.RGBA{
			R: uint8(intensity.Clamp((N.X+1)/2) * 255),
			G: uint8(intensity.Clamp((N.Y+1)/2) * 255),
			B: uint8(intensity.Clamp((N.Z+1)/2) * 255),
			A: 255,
		}
	}

	// Sky
	return color.RGBA{R: 135, G: 206, B: 240, A: 255}
}

type Viewport struct {
	U               *Vec3
	V               *Vec3
	PixelDeltaU     *Vec3
	PixelDeltaV     *Vec3
	UpperLeft       *Vec3
	Pixel00Location *Vec3
}
