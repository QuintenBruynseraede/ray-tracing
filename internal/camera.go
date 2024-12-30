package internal

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"
)

const (
	CAMERA_FOCAL_LENGTH = 1.0
	SAMPLES_PER_PIXEL   = 10
	SAMPLE_SCALE        = 1.0 / SAMPLES_PER_PIXEL
	MAX_DEPTH           = 50
	FOV                 = 100
)

type Camera struct {
	Center      *Vec3
	Viewport    Viewport
	imageWidth  int
	imageHeight int
	maxDepth    int
	fov         int
}

func NewCamera(center *Vec3, imageWidth, imageHeight int) Camera {
	theta := DegToRad(FOV)
	h := math.Tan(theta / 2)
	viewportHeight := 2 * h * CAMERA_FOCAL_LENGTH
	viewPortWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	viewportU := NewVec3(viewPortWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

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
		maxDepth:    MAX_DEPTH,
	}
}

func (c *Camera) Render(image *image.RGBA, world *HittableList) *image.RGBA {
	for y := 0; y < c.imageHeight; y++ {
		for x := 0; x < c.imageWidth; x++ {
			// Use normal ints for intermediate color samples to allow values > 255
			r, g, b := 0, 0, 0
			for sample := 0; sample < SAMPLES_PER_PIXEL; sample++ {
				sampleColor := c.rayColor(c.getRay(x, y), c.maxDepth, world)
				r += int(sampleColor.R)
				g += int(sampleColor.G)
				b += int(sampleColor.B)
			}
			pixelColor := color.RGBA{
				R: uint8(float64(r) * SAMPLE_SCALE),
				G: uint8(float64(g) * SAMPLE_SCALE),
				B: uint8(float64(b) * SAMPLE_SCALE),
			}
			// correctedColor := GammaCorrect(pixelColor)
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

func (c *Camera) rayColor(ray *Ray, depth int, world *HittableList) color.RGBA {
	if depth <= 0 {
		return color.RGBA{}
	}

	hit, hitRecord := world.Hit(ray, &Interval{0.001, math.MaxFloat64})

	if hit {
		scatter, scatterCol, scatterRay := hitRecord.Material.scatter(ray, hitRecord)
		if scatter {
			// Combine material color and scattered ray
			return MultiplyColorValues(scatterCol, c.rayColor(scatterRay, depth-1, world))
		}

		dir := hitRecord.N.Add(RandomUnitVec3()) // Lambertian
		nextBounceColor := c.rayColor(NewRay(hitRecord.P, dir), depth-1, world)
		return MultiplyColor(nextBounceColor, 0.5)
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
