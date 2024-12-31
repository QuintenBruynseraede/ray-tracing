package internal

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

func RandomVec3() *Vec3 {
	return &Vec3{
		X: rand.Float64() * math.MaxFloat64,
		Y: rand.Float64() * math.MaxFloat64,
		Z: rand.Float64() * math.MaxFloat64,
	}
}

func RandomVec3Within(min, max float64) *Vec3 {
	return &Vec3{
		X: min + rand.Float64()*(max-min),
		Y: min + rand.Float64()*(max-min),
		Z: min + rand.Float64()*(max-min),
	}
}

func RandomUnitVec3() *Vec3 {
	for {
		v := RandomVec3Within(-1, 1)
		lensqd := v.LengthSquared()
		if lensqd > 1e-160 && lensqd <= 1 {
			return v.Div(math.Sqrt(lensqd))
		}
	}
}

func RandomOnHemisphere(normal *Vec3) *Vec3 {
	onUnitSphere := RandomUnitVec3()
	if (onUnitSphere.Dot(normal)) > 0.0 { // In same hemisphere as normal
		return onUnitSphere
	} else {
		return onUnitSphere.Mul(-1)
	}
}

func RandomInUnitDisk() *Vec3 {
	for {
		p := NewVec3(-1+rand.Float64()*2, -1+rand.Float64()*2, 0)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

// Returns the reflection of v against a surface with normal n
func Reflect(v *Vec3, n *Vec3) *Vec3 {
	return v.Sub(n.Mul(v.Dot(n)).Mul(2.0))
}

func Refract(v *Vec3, n *Vec3, etaiOverEtat float64) *Vec3 {
	cosTheta := math.Min(v.Mul(-1).Dot(n), 1.0)
	rOutPerp := v.Add(n.Mul(cosTheta)).Mul(etaiOverEtat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (v *Vec3) Add(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v *Vec3) Sub(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v *Vec3) MulVec(u *Vec3) *Vec3 {
	return &Vec3{
		X: v.X * u.X,
		Y: v.Y * u.Y,
		Z: v.Z * u.Z,
	}
}

func (v *Vec3) Mul(scalar float64) *Vec3 {
	return &Vec3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

func (v *Vec3) Div(scalar float64) *Vec3 {
	return &Vec3{
		X: v.X / scalar,
		Y: v.Y / scalar,
		Z: v.Z / scalar,
	}
}

func (v *Vec3) String() string {
	return fmt.Sprintf("v(%v, %v, %v)", v.X, v.Y, v.Z)
}

func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3) Length() float64 {
	return float64(math.Sqrt(v.LengthSquared()))
}

func (v *Vec3) Dot(other *Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v *Vec3) Normalize() *Vec3 {
	length := v.Length()
	if length == 0 {
		return v
	}
	return v.Div(length)
}

func (v *Vec3) NearZero() bool {
	small := 1e-8
	return math.Abs(v.X) < small && math.Abs(v.Y) < small && math.Abs(v.Z) < small
}
