package imgui

import (
	"image/color"

	"github.com/inkyblackness/imgui-go/v2"
)

func Color(c color.RGBA) imgui.Vec4 {
	return imgui.Vec4{
		X: float32(c.R) / 255,
		Y: float32(c.G) / 255,
		Z: float32(c.B) / 255,
		W: float32(c.A) / 255,
	}
}
