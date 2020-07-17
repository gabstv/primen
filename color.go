package primen

import (
	"image/color"

	"github.com/gabstv/primen/internal/z"
)

// ColorFromHex parses a RRGGBBAA (or RRGGBB) hexadecimal color into color.RGBA
//
// Returns color.Transparent on invalid hex inputs.
func ColorFromHex(hexv string) color.RGBA {
	return z.ColorFromHex(hexv)
}
