package z

import (
	"math"
	"math/rand"
)

const rchars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Rs unsafe random string (used for internal IDs)
func Rs() string {
	rs := [16]byte{}
	rand.Read(rs[:])
	rn := [16]rune{}
	lchars := float64(len(rchars) - 1)
	for i := range rs {
		x := (float64(rs[i]) / 255.0) * lchars
		rn[i] = rune(rchars[int(math.Floor(x))])
	}
	return string(rn[:])
}

func S(a string, rest ...string) string {
	if a != "" {
		return a
	}
	for _, v := range rest {
		if v != "" {
			return v
		}
	}
	return ""
}
