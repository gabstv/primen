package renderer

import (
	"strconv"
	"strings"

	"github.com/inkyblackness/imgui-go/v2"
)

func parseWidth(ctx *Context, domval string) (float32, bool) {
	if strings.HasPrefix(domval, "js:") {
		ctx.JS.Set("__column_width", ctx.Stack.MaxWidthD(float32(imgui.ColumnWidth())))
		rawv, err := ctx.JS.RunString(domval[3:])
		if err != nil {
			// TODO: log error to a console
			return 0, false
		}
		vi := rawv.ToNumber().Export()
		if vi == nil {
			return 0, false
		}
		if v, ok := vi.(float64); ok {
			return float32(v), true
		}
		return 0, false
	}
	if strings.HasSuffix(domval, "%") {
		vpct, err := strconv.ParseFloat(domval[:len(domval)-1], 64)
		if err != nil {
			// TODO: log error to a console (?)
			return 0, false
		}
		return float32(vpct/100) * ctx.Stack.MaxWidthD(0), true
	}
	v64, err := strconv.ParseFloat(domval[:len(domval)-1], 64)
	if err != nil {
		// TODO: log error to a console (?)
		return 0, false
	}
	return float32(v64), true
}
