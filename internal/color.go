package internal

type Color = *Vec3

func C(r, g, b float64) Color {
	return NewVec3(r, g, b)
}
