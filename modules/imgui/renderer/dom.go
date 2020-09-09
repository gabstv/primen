package renderer

import (
	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/modules/imgui/style"
	"github.com/inkyblackness/imgui-go/v2"
)

type Context struct {
	Draw  core.DrawCtx
	JS    *goja.Runtime
	Stack *style.Stack
}

func NewContext(ctx core.DrawCtx, jsvm *goja.Runtime) *Context {
	return &Context{
		Draw:  ctx,
		JS:    jsvm,
		Stack: &style.Stack{},
	}
}

type DomRenderFn func(ctx *Context, node dom.ElementNode)

func renderNodeByTagName(ctx *Context, node dom.ElementNode) bool {
	switch node.TagName() {
	case "_root":
		children(ctx, node)
		return true
	case "button":
		button(ctx, node)
		return true
	case "column":
		groupColumn(ctx, node)
		return true
	case "columns":
		groupColumns(ctx, node)
		return true
	case "demowindow":
		demoWindow(ctx, node)
		return true
	case "group":
		group(ctx, node)
		return true
	case "separator":
		imgui.Separator()
		return true
	case "spacing":
		imgui.Spacing()
		return true
	case "window":
		window(ctx, node)
		return true
	}
	return false
}

func Node(ctx *Context, node dom.ElementNode) {
	attrs := node.Attributes()
	// lctx := setNodeLayout(node, data, jsvm) // FIXME: better solution for w="" h=""
	n := pushNodeVariants(ctx, node)
	sn, cn := style.Push(attrs)
	// setNodeLayout(node, data, jsvm)
	defer style.Pop(sn, cn)
	_ = renderNodeByTagName(ctx, node)
	popNodeVariants(ctx, node, n)
}

func group(ctx *Context, node dom.ElementNode) {
	//TODO: common styles and things
	imgui.BeginGroup()
	children(ctx, node)
	imgui.EndGroup()
}

func groupColumns(ctx *Context, node dom.ElementNode) {
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
	children(ctx, node)
}

func groupColumn(ctx *Context, node dom.ElementNode) {
	children(ctx, node)
	//FIXME: get width if defined
	// if lctx.Size.X > 0 {
	// 	imgui.SetColumnWidth(-1, lctx.Size.X)
	// }
	imgui.NextColumn()
}

func pushNodeVariants(ctx *Context, node dom.ElementNode) int {
	attrs := node.Attributes()
	n := 0
	if attrs.HasAttr("width", "w") {
		if v, ok := parseWidth(ctx, attrs.FirstAttr("width", "w")); ok {
			ctx.Stack.PushWidth(v)
			n++
		}
	}
	return n
}

func popNodeVariants(ctx *Context, node dom.ElementNode, n int) {
	for i := 0; i < n; i++ {
		ctx.Stack.Pop()
	}
}

func mf32(v float32, ok bool) float32 { return v }
