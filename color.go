package primen

import (
	"encoding/hex"
	"image/color"
)

// ColorFromHex parses a RRGGBBAA (or RRGGBB) hexadecimal color into color.RGBA
//
// Returns color.Transparent on invalid hex inputs.
func ColorFromHex(hexv string) color.RGBA {
	if hexv == "" {
		return color.RGBA{}
	}
	if hexv[0] == '#' {
		hexv = hexv[1:]
	}
	if len(hexv) != 8 && len(hexv) != 6 {
		return color.RGBA{}
	}
	r, _ := hex.DecodeString(hexv[:2])
	if len(r) != 1 {
		return color.RGBA{}
	}
	g, _ := hex.DecodeString(hexv[2:4])
	if len(g) != 1 {
		return color.RGBA{}
	}
	b, _ := hex.DecodeString(hexv[4:6])
	if len(b) != 1 {
		return color.RGBA{}
	}
	a := []byte{255}
	if len(hexv) == 8 {
		a, _ = hex.DecodeString(hexv[6:])
		if len(a) != 1 {
			return color.RGBA{}
		}
	}
	return color.RGBA{r[0], g[0], b[0], a[0]}
}
