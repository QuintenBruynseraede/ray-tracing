package scenes

import (
	"math/rand/v2"

	"github.com/quintenbruynseraede/ray-tracing/internal"
	. "github.com/quintenbruynseraede/ray-tracing/internal"
)

func LoadDefocusScene() (Camera, *HittableList) {
	groundMat := Lambertian{C(0.7, 0.7, 0.7)}
	ground := NewSphere(internal.NewVec3(0, -1000, 0), 1000, groundMat)

	world := NewHittableList(ground)

	N := 10
	for a := 0; a < N; a++ {
		chooseMat := rand.Float64()
		center := NewVec3(0, 1, -10+3*float64(a))

		if chooseMat < 0.25 {
			// Diffuse
			world.Add(NewSphere(center, 1, Lambertian{Albedo: RandomColor()}))
			world.Add(NewSphere(center, 1, Lambertian{Albedo: RandomColor()}))
		} else {
			// Metal
			world.Add(NewSphere(center, 1, Metal{Albedo: RandomColor(), Fuzz: rand.Float64() / 4}))
		}
	}

	camera := NewCamera(
		1200,
		675,
		50,                  // FOV
		NewVec3(-8, 8, -15), // Look from
		NewVec3(0, -10, 15),
		NewVec3(0, 1, 0),
		1.0,
		20.0, // Focal dist
	)

	return camera, world
}
