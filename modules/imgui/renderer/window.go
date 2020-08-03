package renderer

import (
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
	"github.com/gabstv/primen/modules/imgui/common"
	"github.com/gabstv/primen/modules/imgui/style"
	"github.com/inkyblackness/imgui-go/v2"
)

// Window handle imgui's Begin()
func Window(ctx *Context, node dom.ElementNode) {
	attrs := node.Attributes()
	if node.ID() == "" {
		println("warning: window didn't have an ID")
		node.SetAttribute("id", z.Rs())
	}
	wname := attrs["name"]
	if wname == "" {
		wname = node.ID()
	}

	show := node.Attributes().BoolD("visible", true)
	lshow := show
	if show {
		style.SetupWindowPos(ctx.Draw, attrs, ctx.JS)
		imgui.BeginV(wname, &show, parseWindowFlags(attrs))
		Children(ctx, node)
		imgui.End()
	}
	if lshow != show {
		// the window was closed by the X button
		node.SetAttribute("visible", common.TernaryString(show, "true", "false"))
	}
}

func DemoWindow(ctx *Context, node dom.ElementNode) {
	if node.ID() == "" {
		println("warning: demo window didn't have an ID")
		node.SetAttribute("id", z.Rs())
	}
	show := node.Attributes().BoolD("visible", true)
	lshow := show
	if show {
		imgui.ShowDemoWindow(&show)
	}
	if lshow != show {
		// the window was closed by the X button
		node.SetAttribute("visible", common.TernaryString(show, "true", "false"))
	}
}
