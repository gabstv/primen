package z

import "strconv"

func Float64(v string, defaultf float64) float64 {
	if v == "" {
		return defaultf
	}
	if vf, err := strconv.ParseFloat(v, 64); err == nil {
		return vf
	}
	return defaultf
}

func Float32(v string, defaultf float32) float32 {
	if v == "" {
		return defaultf
	}
	if vf, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(vf)
	}
	return defaultf
}

func Float32V(v string) (float32, bool) {
	if v == "" {
		return 0, false
	}
	if vf, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(vf), true
	}
	return 0, false
}
