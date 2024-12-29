package internal

import "math"

type Interval struct {
	Min float64
	Max float64
}

func (i *Interval) Size() float64 {
	return i.Max - i.Min
}

func (i *Interval) Contains(x float64) bool {
	return x >= i.Min && x <= i.Max
}

func (i *Interval) Surrounds(x float64) bool {
	return x > i.Min && x < i.Max
}

func (i *Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	} else if x > i.Max {
		return i.Max
	}
	return x
}

var (
	Empty    = Interval{math.MaxFloat64, -math.MaxFloat64}
	Universe = Interval{-math.MaxFloat64, math.MaxFloat64}
)
