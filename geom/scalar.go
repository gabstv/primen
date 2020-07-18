package geom

import "math"

func ScalarEqualsEpsilon(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	if d < math.SmallestNonzeroFloat64 {
		return true
	}
	return d < epsilon
}
