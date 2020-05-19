package core

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
