package core

import "image/color"

func applyOrigin(length, origin float64) float64 {
	return length * (-origin)
}

// nonzeroval will return the first non zero value.
// It will return 0 if all input values are 0
func nonzeroval(vals ...float64) float64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

// Clamp returns x clamped to the interval [min, max].
//
// If x is less than min, min is returned. If x is more than max, max is returned. Otherwise, x is
// returned.
func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func Scalef(f, scale float64) float64 {
	return f * scale
}

func Scaleb(b uint8, scale float64) uint8 {
	return uint8(float64(b) * scale)
}

func Lerpf(a, b, t float64) float64 {
	return Scalef(a, 1-t) + Scalef(b, t)
}

func Lerpc(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: Scaleb(a.R, 1-t) + Scaleb(b.R, t),
		G: Scaleb(a.G, 1-t) + Scaleb(b.G, t),
		B: Scaleb(a.B, 1-t) + Scaleb(b.B, t),
		A: Scaleb(a.A, 1-t) + Scaleb(b.A, t),
	}
}

func intmax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
