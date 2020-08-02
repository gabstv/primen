package imgui

import (
	"errors"

	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/js"
	"github.com/gabstv/primen/dom"
	domrenderer "github.com/gabstv/primen/modules/imgui/renderer"
)

type UI struct {
	id       UID
	document dom.ElementNode
	jsvm     *goja.Runtime
}

func newUI(id UID, doc []dom.Node) *UI {
	//goja.New()
	ui := &UI{
		id:       id,
		document: dom.Element("_root", nil, doc...),
		jsvm:     goja.New(),
	}
	ui.inlineJs()
	return ui
}

func (ui *UI) inlineJs() {
	vm := ui.jsvm
	//
	js.Stdlib(lastEngine, vm)

	vm.Set("dispatch", func(name goja.Value, value goja.Value) error {
		if vs, ok := name.Export().(string); ok {
			lastEngine.DispatchEvent(vs, value.Export())
			return nil
		}
		return errors.New("event name must be a string")
	})
	vm.Set("show", func(id string) {
		if n := ui.document.FindChildByID(id); n != nil {
			if n.Type() == dom.NodeElement {
				n.(dom.ElementNode).SetAttribute("visible", "true")
			}
		}
	})
	vm.Set("hide", func(id string) {
		if n := ui.document.FindChildByID(id); n != nil {
			if n.Type() == dom.NodeElement {
				n.(dom.ElementNode).SetAttribute("visible", "false")
			}
		}
	})
	vm.Set("setattr", func(id, name, value string) {
		if n := ui.document.FindChildByID(id); n != nil {
			if n.Type() == dom.NodeElement {
				n.(dom.ElementNode).SetAttribute(name, value)
				//TODO: force imgui to change position if attr is x, y, width, height
			}
		}
	})
}

func (ui *UI) Render(ctx core.DrawCtx) {
	w, h := ctx.Renderer().Screen().Size()
	ui.jsvm.Set("__width", w)
	ui.jsvm.Set("__height", h)
	// renderNode(ctx, ui.document, ui.data, ui.jsvm)
	rctx := domrenderer.NewContext(ctx, ui.jsvm)
	domrenderer.Node(rctx, ui.document)
}

// func renderNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime) {
// 	attrs := node.Attributes()
// 	lctx := setNodeLayout(node, data, jsvm)
// 	sn, cn := style.Push(attrs, data)
// 	setNodeLayout(node, data, jsvm)
// 	defer style.Pop(sn, cn)
// 	switch node.TagName() {
// 	case "_root":
// 		renderRootNode(ctx, node, data, jsvm, lctx)
// 	case "window":
// 		renderWindowNode(ctx, node, data, jsvm, lctx)
// 	case "demowindow":
// 		renderDemoWindowNode(ctx, node, data, jsvm, lctx)
// 	case "group":
// 		renderGroupNode(ctx, node, data, jsvm, lctx)
// 	case "columns":
// 		renderGroupColumns(ctx, node, data, jsvm, lctx)
// 	case "column":
// 		renderGroupColumn(ctx, node, data, jsvm, lctx)
// 	case "separator":
// 		imgui.Separator()
// 	case "spacing":
// 		imgui.Spacing()
// 	case "button":
// 		renderButtonNode(ctx, node, data, jsvm, lctx)
// 	}
// }

// func renderButtonNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	attrs := node.Attributes()
// 	if imgui.ButtonV(z.S(attrs["label"], node.FirstChildAsText()), lctx.Size) {
// 		bubbles := true
// 		if jsv := attrs["onclick"]; jsv != "" {
// 			rv, err := jsvm.RunString(jsv)
// 			if err != nil {
// 				println("js err: " + err.Error())
// 				bubbles = false
// 			} else {
// 				println("js fn ok")
// 				if rv == nil || !rv.ToBoolean() {
// 					bubbles = false
// 				}
// 			}
// 		}
// 		if bubbles {
// 			println("bubbles: button click!!!")
// 		}
// 	}
// }

// func renderRootNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	renderChildren(ctx, node, data, jsvm)
// }

// func renderWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	attrs := node.Attributes()
// 	if node.ID() == "" {
// 		println("warning: window didn't have an ID")
// 		node.SetAttribute("id", z.Rs())
// 	}
// 	wname := attrs["name"]
// 	if wname == "" {
// 		wname = node.ID()
// 	}
// 	show := data.MustBool(common.VarVisible(node.ID()), node.Attributes().BoolD("visible", true))
// 	lshow := show
// 	if show {
// 		style.SetupWindowPos(ctx, attrs, data, jsvm)
// 		imgui.BeginV(wname, &show, parseWindowFlags(attrs))
// 		renderChildren(ctx, node, data, jsvm)
// 		imgui.End()
// 	}
// 	if lshow != show {
// 		// the window was closed by the X button
// 		data.Set(common.VarVisible(node.ID()), show)
// 	}
// }

// func renderGroupNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	imgui.BeginGroup()
// 	renderChildren(ctx, node, data, jsvm)
// 	imgui.EndGroup()
// }

// func renderGroupColumns(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	attr := node.Attributes()
// 	if v := attr.IntD("count", 0); v > 0 {
// 		imgui.ColumnsV(v, attr.String("label"), attr.BoolD("border", false))
// 	} else {
// 		nc := 0
// 		// count columns (slower)
// 		for _, c := range node.Children() {
// 			if c.Type() == dom.NodeElement {
// 				if c.(dom.ElementNode).TagName() == "column" {
// 					nc++
// 				}
// 			}
// 		}
// 		imgui.ColumnsV(nc, attr.String("label"), attr.BoolD("border", false))
// 	}
// 	renderChildren(ctx, node, data, jsvm)
// }

// func renderGroupColumn(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	renderChildren(ctx, node, data, jsvm)
// 	if lctx.Size.X > 0 {
// 		imgui.SetColumnWidth(-1, lctx.Size.X)
// 	}
// 	imgui.NextColumn()
// }

// func renderChildren(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime) {
// 	for _, child := range node.Children() {
// 		if child.Type() == dom.NodeElement {
// 			renderNode(ctx, child.(dom.ElementNode), data, jsvm)
// 		} else if child.Type() == dom.NodeText {
// 			imgui.Text(child.(dom.TextNode).Text())
// 		}
// 	}
// }

// func renderDemoWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *store.DB, jsvm *goja.Runtime, lctx layoutContext) {
// 	if node.ID() == "" {
// 		println("warning: window didn't have an ID")
// 		node.SetAttribute("id", z.Rs())
// 	}
// 	show := data.MustBool(common.VarVisible(node.ID()), node.Attributes().BoolD("visible", true))
// 	lshow := show
// 	if show {
// 		imgui.ShowDemoWindow(&show)
// 	}
// 	if lshow != show {
// 		// the window was closed by the X button
// 		data.Set(common.VarVisible(node.ID()), show)
// 	}
// }
