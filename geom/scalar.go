package geom

import "math"

const (
	// Epsilon is the default value to compare is scalars are nearly equal
	// It can be as low as math.SmallestNonzeroFloat64
	Epsilon = 1e-16
)

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
