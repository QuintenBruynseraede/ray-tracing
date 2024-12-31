package scenes

import (
	"math"
	"math/rand/v2"

	"github.com/quintenbruynseraede/ray-tracing/internal"
	. "github.com/quintenbruynseraede/ray-tracing/internal"
)

func LoadCircleScene() (Camera, *HittableList) {
	groundMat := Lambertian{C(0.5, 0.5, 0.5)}
	ground := NewSphere(internal.NewVec3(0, -1000, 0), 1000, groundMat)

	world := NewHittableList(ground)

	N := 14
	R := 10.0
	angleStep := 2 * math.Pi / float64(N)
	for a := 0; a < N; a++ {
		angle := angleStep * float64(a)
		chooseMat := rand.Float64()
		center := NewVec3(R*math.Cos(angle), 1, R*math.Sin(angle))

		if chooseMat < 0.25 {
			// Diffuse
			world.Add(NewSphere(center, 1, Lambertian{Albedo: RandomColor()}))
			world.Add(NewSphere(center, 1, Lambertian{Albedo: RandomColor()}))
		} else if chooseMat < 0.75 {
			// Metal
			world.Add(NewSphere(center, 1, Metal{Albedo: RandomColor(), Fuzz: rand.Float64() / 4}))
		} else {
			// Glass
			world.Add(NewSphere(center, 1, Dielectric{RefractionIndex: 1.5}))
		}
	}

	// Central circle
	mat3 := Metal{Albedo: C(0.9, 0.9, 0.9), Fuzz: 0.0}
	world.Add(NewSphere(NewVec3(0, 6, 0), 5.0, mat3))

	camera := NewCamera(
		1200,
		675,
		70,
		NewVec3(18, 6.5, 16),
		NewVec3(7, -3, 0),
		NewVec3(0, 1, 0),
		0.0,
		10.0,
	)

	return camera, world
}
