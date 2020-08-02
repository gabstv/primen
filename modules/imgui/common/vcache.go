package common

import (
	"strconv"
	"strings"

	"github.com/gabstv/primen/internal/z"
	"github.com/inkyblackness/imgui-go/v2"
)

func ParseVec2(v string) (imgui.Vec2, bool) {
	if v == "" {
		return imgui.Vec2{}, false
	}
	i := strings.IndexByte(v, ',')
	if i <= 0 {
		f64, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return imgui.Vec2{}, false
		}
		out := imgui.Vec2{
			X: float32(f64),
			Y: float32(f64),
		}
		return out, true
	}
	x, err1 := strconv.ParseFloat(strings.TrimSpace(v[:i]), 32)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(v[i+1:]), 32)
	if err1 != nil || err2 != nil {
		return imgui.Vec2{}, false
	}
	out := imgui.Vec2{
		X: float32(x),
		Y: float32(y),
	}
	return out, true
}

func ParseVec4(v string) (imgui.Vec4, bool) {
	if v == "" {
		return imgui.Vec4{}, false
	}
	if v[0] == '#' {
		rgba := z.ColorFromHex(v)
		out := imgui.Vec4{
			X: float32(rgba.R) / 255,
			Y: float32(rgba.G) / 255,
			Z: float32(rgba.B) / 255,
			W: float32(rgba.A) / 255,
		}
		return out, true
	}
	slcs := strings.Split(v, ",")
	if len(slcs) == 3 {
		f64x, _ := strconv.ParseFloat(strings.TrimSpace(slcs[0]), 32)
		f64y, _ := strconv.ParseFloat(strings.TrimSpace(slcs[1]), 32)
		f64z, _ := strconv.ParseFloat(strings.TrimSpace(slcs[2]), 32)
		out := imgui.Vec4{
			X: float32(f64x),
			Y: float32(f64y),
			Z: float32(f64z),
			W: 1.0,
		}
		return out, true
	} else if len(slcs) == 4 {
		f64x, _ := strconv.ParseFloat(strings.TrimSpace(slcs[0]), 32)
		f64y, _ := strconv.ParseFloat(strings.TrimSpace(slcs[1]), 32)
		f64z, _ := strconv.ParseFloat(strings.TrimSpace(slcs[2]), 32)
		f64w, _ := strconv.ParseFloat(strings.TrimSpace(slcs[3]), 32)
		out := imgui.Vec4{
			X: float32(f64x),
			Y: float32(f64y),
			Z: float32(f64z),
			W: float32(f64w),
		}
		return out, true
	}
	return imgui.Vec4{}, false
}

func VarVisible(nodeid string) string {
	return nodeid + "_visible"
}

func VarPosition(nodeid string) string {
	return nodeid + "_position"
}
