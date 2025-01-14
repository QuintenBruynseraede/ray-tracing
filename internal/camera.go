package internal

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"sync"
)

const (
	SAMPLES_PER_PIXEL = 500
	SAMPLE_SCALE      = 1.0 / SAMPLES_PER_PIXEL
	MAX_DEPTH         = 500
)

type Camera struct {
	Center       *Vec3
	u, v, w      *Vec3 // Camera frame basis vectors
	defocusDiskU *Vec3
	defocusDiskV *Vec3
	Viewport     Viewport
	imageWidth   int
	imageHeight  int
	maxDepth     int
	fov          int
	defocusAngle float64
}

func NewCamera(imageWidth, imageHeight int, fov int, lookFrom, lookAt, vup *Vec3, defocusAngle, focusDist float64) Camera {
	// Basis vectors
	w := (lookFrom.Sub(lookAt)).Normalize()
	u := vup.Cross(w).Normalize()
	v := w.Cross(u)

	// FOV
	theta := DegToRad(float64(fov))
	h := math.Tan(theta / 2)
	viewportHeight := 2 * h * focusDist
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	viewportU := u.Mul(viewportWidth)
	viewportV := v.Mul(-1.0 * viewportHeight)

	// Calculate horizontal and vertical delta vectors from pixel to pixel
	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	// Calculate the location of the upper left pixel
	viewportUpperLeft := (lookFrom.
		Sub(w.Mul(focusDist)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2)))

	// Calculate the location of the lower right pixel
	pixel00Location := viewportUpperLeft.
		Add(pixelDeltaU.Mul(0.5)).
		Add(pixelDeltaV.Mul(0.5))

	// Defocus basis vectors
	defocusRadius := focusDist * math.Tan(DegToRad(defocusAngle/2))

	viewport := Viewport{
		U:               viewportU,
		V:               viewportV,
		PixelDeltaU:     pixelDeltaU,
		PixelDeltaV:     pixelDeltaV,
		UpperLeft:       viewportUpperLeft,
		Pixel00Location: pixel00Location,
	}

	return Camera{
		Center:       lookFrom,
		Viewport:     viewport,
		imageWidth:   imageWidth,
		imageHeight:  imageHeight,
		maxDepth:     MAX_DEPTH,
		w:            w,
		u:            u,
		v:            v,
		defocusDiskU: u.Mul(defocusRadius),
		defocusDiskV: v.Mul(defocusRadius),
		defocusAngle: defocusAngle,
	}
}

func (c *Camera) Render(image *image.RGBA, world *HittableList) *image.RGBA {
	imArray := make([][]Color, c.imageHeight)
	for i := range imArray {
		imArray[i] = make([]Color, c.imageWidth)
	}

	var wg sync.WaitGroup

	for y := 0; y < c.imageHeight; y++ {
		for x := 0; x < c.imageWidth; x++ {
			wg.Add(1)
			// One goroutine per pixel
			go func(x, y int) {
				defer wg.Done()

				// Use normal ints for intermediate color samples to allow values > 255
				// TODO: just use vector math here
				r, g, b := 0.0, 0.0, 0.0
				for sample := 0; sample < SAMPLES_PER_PIXEL; sample++ {
					sampleColor := c.rayColor(c.getRay(x, y), c.maxDepth, world)
					r += sampleColor.X
					g += sampleColor.Y
					b += sampleColor.Z
				}
				imArray[y][x] = C(r*SAMPLE_SCALE, g*SAMPLE_SCALE, b*SAMPLE_SCALE).Mul(255.0)
			}(x, y)
		}
	}

	wg.Wait()

	for y := 0; y < c.imageHeight; y++ {
		for x := 0; x < c.imageWidth; x++ {
			image.Set(x, y, color.RGBA{
				R: uint8(imArray[y][x].X),
				G: uint8(imArray[y][x].Y),
				B: uint8(imArray[y][x].Z),
			})
		}
	}

	return image
}

func (c *Camera) getRay(x, y int) *Ray {
	offsetX, offsetY := sampleSquare()
	pixelSample := c.Viewport.Pixel00Location.
		Add(c.Viewport.PixelDeltaU.Mul(float64(x) + offsetX)).
		Add(c.Viewport.PixelDeltaV.Mul(float64(y) + offsetY))

	var rayOrigin *Vec3
	if c.defocusAngle <= 0 {
		rayOrigin = c.Center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	return NewRay(rayOrigin, pixelSample.Sub(rayOrigin))
}

func sampleSquare() (float64, float64) {
	return rand.Float64() - 0.5, rand.Float64() - 0.5
}

func (c *Camera) rayColor(ray *Ray, depth int, world *HittableList) Color {
	if depth <= 0 {
		return C(0, 0, 0)
	}

	hit, hitRecord := world.Hit(ray, &Interval{0.001, math.MaxFloat64})

	if hit {
		scatter, scatterCol, scatterRay := hitRecord.Material.scatter(ray, hitRecord)
		if scatter {
			// Combine material color and scattered ray
			return scatterCol.MulVec(c.rayColor(scatterRay, depth-1, world))
		}

		var origin *Vec3
		if c.defocusAngle <= 0 {
			origin = c.Center
		} else {
			origin = c.defocusDiskSample()
		}
		dir := hitRecord.N.Add(RandomUnitVec3()) // Lambertian
		nextBounceColor := c.rayColor(NewRay(origin, dir), depth-1, world)
		return nextBounceColor.Mul(0.5)
	}

	// Sky
	return C(0.53, 0.81, 0.94)
}

// defocusDiskSample returns a random point in the camera defocus disk
func (c *Camera) defocusDiskSample() *Vec3 {
	p := RandomInUnitDisk()
	return c.Center.
		Add(c.defocusDiskU.Mul(p.X)).
		Add(c.defocusDiskV.Mul(p.Y))
}

type Viewport struct {
	U               *Vec3
	V               *Vec3
	PixelDeltaU     *Vec3
	PixelDeltaV     *Vec3
	UpperLeft       *Vec3
	Pixel00Location *Vec3
	FocalLength     float64
}
