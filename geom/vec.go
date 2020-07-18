package geom

import "math"

const (
	// Epsilon is the default value to compare is scalars are nearly equal
	// It can be as low as math.SmallestNonzeroFloat64
	Epsilon = 0.0000000000000001
)

// ZV = Vec{0,0}
var ZV = Vec{}

// Vec is a 2d Vector
type Vec struct {
	X float64
	Y float64
}

// Equals returns v == other
func (v Vec) Equals(other Vec) bool {
	return v.X == other.X && v.Y == other.Y
}

func (v Vec) EqualsEpsilon(other Vec) bool {
	return v.EqualsEpsilon2(other, Epsilon)
}

func (v Vec) EqualsEpsilon2(other Vec, epsilon float64) bool {
	if v.Equals(other) {
		return true
	}
	return ScalarEqualsEpsilon(v.X, other.X, epsilon) &&
		ScalarEqualsEpsilon(v.Y, other.Y, epsilon)
}

// IsZero returns true if both axes are 0
func (v Vec) IsZero() bool {
	return v.Equals(ZV)
}

// Dot product
func (v Vec) Dot(other Vec) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross product
func (v Vec) Cross(other Vec) float64 {
	return v.X*other.Y - other.X*v.Y
}

// Magnitude = length
func (v Vec) Magnitude() float64 {
	return math.Hypot(v.X, v.Y)
}

// Normalized normalizes a vector.
// Also known as direction, unit vector.
func (v Vec) Normalized() Vec {
	if v.X == 0 && v.Y == 0 {

	}
	m := v.Magnitude()
	return Vec{v.X / m, v.Y / m}
}

// Scaled returns {v.X * s, v.Y * s}
func (v Vec) Scaled(s float64) Vec {
	return Vec{v.X * s, v.Y * s}
}

func (v Vec) Applyed() {

}
