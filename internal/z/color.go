package z

import (
	"encoding/hex"
	"image/color"
)

var White = color.RGBA{255, 255, 255, 255}
var Black = color.RGBA{0, 0, 0, 255}

func Color(hexv string, defaultc color.RGBA) color.RGBA {
	if hexv == "" {
		return defaultc
	}
	c, ok := colorfhex(hexv)
	if !ok {
		return defaultc
	}
	return c
}

// ColorFromHex parses a RRGGBBAA (or RRGGBB) hexadecimal color into color.RGBA
//
// Returns color.Transparent on invalid hex inputs.
func ColorFromHex(hexv string) color.RGBA {
	c, _ := colorfhex(hexv)
	return c
}

func colorfhex(hexv string) (color.RGBA, bool) {
	if hexv == "" {
		return color.RGBA{}, false
	}
	if hexv[0] == '#' {
		hexv = hexv[1:]
	}
	if len(hexv) != 8 && len(hexv) != 6 {
		return color.RGBA{}, false
	}
	r, _ := hex.DecodeString(hexv[:2])
	if len(r) != 1 {
		return color.RGBA{}, false
	}
	g, _ := hex.DecodeString(hexv[2:4])
	if len(g) != 1 {
		return color.RGBA{}, false
	}
	b, _ := hex.DecodeString(hexv[4:6])
	if len(b) != 1 {
		return color.RGBA{}, false
	}
	a := []byte{255}
	if len(hexv) == 8 {
		a, _ = hex.DecodeString(hexv[6:])
		if len(a) != 1 {
			return color.RGBA{}, false
		}
	}
	return color.RGBA{r[0], g[0], b[0], a[0]}, true
}
