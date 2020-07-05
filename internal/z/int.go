package z

import "strconv"

func Int(v string, defaulti int) int {
	if v == "" {
		return defaulti
	}
	if vi, err := strconv.Atoi(v); err == nil {
		return vi
	}
	return defaulti
}
