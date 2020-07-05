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
