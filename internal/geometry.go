package internal

import "math"

type HitRecord struct {
	P         *Vec3
	N         *Vec3
	T         float64
	FrontFace bool
}

type Hittable interface {
	Hit(ray *Ray, rayT *Interval) (bool, HitRecord)
}

type Sphere struct {
	Center *Vec3
	Radius float64
}

func NewSphere(center *Vec3, radius float64) *Sphere {
	return &Sphere{
		Center: center,
		Radius: radius,
	}
}

func (s *Sphere) Hit(ray *Ray, rayT *Interval) (bool, HitRecord) {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return false, HitRecord{}
	} else {
		sqrtd := math.Sqrt(discriminant)
		root := (h - sqrtd) / a
		if !rayT.Surrounds(root) {
			root = (h + sqrtd) / a
			if !rayT.Surrounds(root) {
				return false, HitRecord{}
			}
		}

		T := root
		P := ray.At(T)
		N := P.Sub(s.Center).Div(s.Radius)
		outwardNormal := P.Sub(s.Center).Div(s.Radius)

		hit := HitRecord{
			P: P,
			T: T,
			N: N,
		}
		hit.SetFaceNormal(ray, outwardNormal)
		return true, hit
	}
}

func (r *HitRecord) SetFaceNormal(ray *Ray, outwardNormal *Vec3) {
	r.FrontFace = ray.Direction.Dot(outwardNormal) < 0
	if r.FrontFace {
		r.N = outwardNormal
	} else {
		r.N = outwardNormal.Mul(-1)
	}
}

type HittableList struct {
	Objects []Hittable
}

func NewHittableList(objects ...Hittable) *HittableList {
	return &HittableList{
		Objects: objects,
	}
}

func (h *HittableList) Hit(ray *Ray, rayT *Interval) (bool, HitRecord) {
	closestHit := HitRecord{}
	closestHit.T = rayT.Max
	hasHit := false

	for _, object := range h.Objects {
		hit, hitRecord := object.Hit(ray, &Interval{rayT.Min, closestHit.T})
		if hit {
			closestHit = hitRecord
			hasHit = true
		}
	}
	return hasHit, closestHit
}

func (h *HittableList) Add(object Hittable) {
	h.Objects = append(h.Objects, object)
}
