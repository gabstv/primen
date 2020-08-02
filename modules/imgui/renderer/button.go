package renderer

import (
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
	"github.com/inkyblackness/imgui-go/v2"
)

func Button(ctx *Context, node dom.ElementNode) {
	attrs := node.Attributes()
	//FIXME: get size if defined width and height
	tempsize := imgui.Vec2{}
	if imgui.ButtonV(z.S(attrs.String("label"), node.FirstChildAsText()), tempsize) {
		bubbles := true
		if jsv := attrs["onclick"]; jsv != "" {
			rv, err := ctx.JS.RunString(jsv)
			if err != nil {
				println("js err: " + err.Error())
				bubbles = false
			} else {
				println("js fn ok")
				if rv == nil || !rv.ToBoolean() {
					bubbles = false
				}
			}
		}
		if bubbles {
			println("bubbles: button click!!!")
		}
	}
}
