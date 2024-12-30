package internal

import (
	"image/color"
	"math"
	"math/rand/v2"
)

type Material interface {
	scatter(rayIn *Ray, hit HitRecord) (bool, color.RGBA, *Ray)
}

type Lambertian struct {
	Albedo color.RGBA
}

func (l Lambertian) scatter(rayIn *Ray, hit HitRecord) (bool, color.RGBA, *Ray) {
	scatterDirection := hit.N.Add(RandomUnitVec3())
	if scatterDirection.NearZero() {
		scatterDirection = hit.N // Prevent edge cases later on
	}
	return true, l.Albedo, NewRay(hit.P, scatterDirection)
}

type Metal struct {
	Albedo color.RGBA
	Fuzz   float64
}

func (m Metal) scatter(rayIn *Ray, hit HitRecord) (bool, color.RGBA, *Ray) {
	reflected := Reflect(rayIn.Direction, hit.N)
	reflected = reflected.Normalize().Add(RandomUnitVec3().Mul(m.Fuzz))

	scattered := NewRay(hit.P, reflected)
	if scattered.Direction.Dot(hit.N) > 0 {
		return true, m.Albedo, scattered
	} else {
		return false, color.RGBA{}, nil
	}
}

type Dielectric struct {
	RefractionIndex float64
}

func (d Dielectric) scatter(rayIn *Ray, hit HitRecord) (bool, color.RGBA, *Ray) {
	var ri float64
	if hit.FrontFace {
		ri = 1.0 / d.RefractionIndex
	} else {
		ri = d.RefractionIndex
	}

	unitDirection := rayIn.Direction.Normalize()
	cosTheta := math.Min(unitDirection.Mul(-1).Dot(hit.N), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1.0
	var direction *Vec3
	if cannotRefract || d.reflectance(cosTheta, ri) > rand.Float64() {
		direction = Reflect(unitDirection, hit.N)
	} else {
		direction = Refract(unitDirection, hit.N, ri)
	}

	return true, color.RGBA{255, 255, 255, 255}, NewRay(hit.P, direction)
}

func (d Dielectric) reflectance(cos float64, ri float64) float64 {
	// Using Schlick's approximation
	r0 := (1 - ri) / (1 + ri)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cos), 5)
}
