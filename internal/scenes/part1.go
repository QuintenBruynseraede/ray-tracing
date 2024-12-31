package scenes

import (
	"math/rand/v2"

	"github.com/quintenbruynseraede/ray-tracing/internal"
	. "github.com/quintenbruynseraede/ray-tracing/internal"
)

func LoadPart1FinalRender() (Camera, *HittableList) {
	groundMat := Lambertian{C(0.5, 0.5, 0.5)}
	ground := NewSphere(internal.NewVec3(0, -1000, 0), 1000, groundMat)

	world := NewHittableList(ground)

	for a := 0; a < 11; a++ {
		for b := 0; b < 11; b++ {
			chooseMat := rand.Float64()
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Length() > 0.9 {
				if chooseMat < 0.8 {
					// Diffuse
					world.Add(NewSphere(center, 0.2, Lambertian{Albedo: RandomColor()}))
					world.Add(NewSphere(center, 0.2, Lambertian{Albedo: RandomColor()}))
				} else if chooseMat < 0.95 {
					// Metal
					world.Add(NewSphere(center, 0.2, Metal{Albedo: RandomColor(), Fuzz: rand.Float64() / 2}))
				} else {
					// Glass
					world.Add(NewSphere(center, 0.2, Dielectric{RefractionIndex: 1.5}))
				}
			}
		}
	}

	mat1 := Dielectric{RefractionIndex: 1.5}
	mat2 := Lambertian{Albedo: C(0.4, 0.2, 0.1)}
	mat3 := Metal{Albedo: C(0.7, 0.6, 0.5), Fuzz: 0.0}

	world.Add(NewSphere(NewVec3(0, 1, 0), 1.0, mat1))
	world.Add(NewSphere(NewVec3(-4, 1, 0), 1.0, mat2))
	world.Add(NewSphere(NewVec3(4, 1, 0), 1.0, mat3))

	camera := NewCamera(
		1200,
		675,
		50,
		NewVec3(13, 1.5, 3),
		NewVec3(0, -2, 0),
		NewVec3(0, 1, 0),
		0.6,
		10.0,
	)

	return camera, world
}

func RandomColor() Color {
	return C(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64())
}
