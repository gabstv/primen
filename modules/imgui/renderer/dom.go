package renderer

import (
	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/modules/imgui/style"
	"github.com/inkyblackness/imgui-go/v2"
)

type Context struct {
	Draw core.DrawCtx
	JS   *goja.Runtime
}

func NewContext(ctx core.DrawCtx, jsvm *goja.Runtime) *Context {
	return &Context{
		Draw: ctx,
		JS:   jsvm,
	}
}

type DomRenderFn func(ctx *Context, node dom.ElementNode)

func renderNodeByTagName(ctx *Context, node dom.ElementNode) bool {
	switch node.TagName() {
	case "_root":
		Children(ctx, node)
		return true
	case "button":
		Button(ctx, node)
		return true
	case "column":
		GroupColumn(ctx, node)
		return true
	case "columns":
		GroupColumns(ctx, node)
		return true
	case "demowindow":
		DemoWindow(ctx, node)
		return true
	case "group":
		Group(ctx, node)
		return true
	case "separator":
		imgui.Separator()
		return true
	case "spacing":
		imgui.Spacing()
		return true
	case "window":
		Window(ctx, node)
		return true
	}
	return false
}

func Node(ctx *Context, node dom.ElementNode) {
	attrs := node.Attributes()
	// lctx := setNodeLayout(node, data, jsvm) // FIXME: better solution for w="" h=""
	sn, cn := style.Push(attrs)
	// setNodeLayout(node, data, jsvm)
	defer style.Pop(sn, cn)
	_ = renderNodeByTagName(ctx, node)
}

func Group(ctx *Context, node dom.ElementNode) {
	//TODO: common styles and things
	imgui.BeginGroup()
	Children(ctx, node)
	imgui.EndGroup()
}

func GroupColumns(ctx *Context, node dom.ElementNode) {
	attr := node.Attributes()
	if v := attr.IntD("count", 0); v > 0 {
		imgui.ColumnsV(v, attr.String("label"), attr.BoolD("border", false))
	} else {
		nc := 0
		// count columns (slower)
		for _, c := range node.Children() {
			if c.Type() == dom.NodeElement {
				if c.(dom.ElementNode).TagName() == "column" {
					nc++
				}
			}
		}
		imgui.ColumnsV(nc, attr.String("label"), attr.BoolD("border", false))
	}
	Children(ctx, node)
}

func GroupColumn(ctx *Context, node dom.ElementNode) {
	Children(ctx, node)
	//FIXME: get width if defined
	// if lctx.Size.X > 0 {
	// 	imgui.SetColumnWidth(-1, lctx.Size.X)
	// }
	imgui.NextColumn()
}
