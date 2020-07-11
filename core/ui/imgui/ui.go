package imgui

import (
	"errors"
	"strings"

	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
	"github.com/inkyblackness/imgui-go/v2"
)

type UI struct {
	id       UID
	data     *uiMemory
	document dom.ElementNode
	jsvm     *goja.Runtime
}

func newUI(id UID, doc dom.ElementNode) *UI {
	//goja.New()
	ui := &UI{
		id:       id,
		document: doc,
		data:     newUIMemory(),
		jsvm:     goja.New(),
	}
	ui.inlineJs()
	return ui
}

func (ui *UI) inlineJs() {
	vm := ui.jsvm
	//
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
				ui.data.Set(varVisible(n.(dom.ElementNode).ID()), true)
			}
		}
	})
	vm.Set("hide", func(id string) {
		if n := ui.document.FindChildByID(id); n != nil {
			if n.Type() == dom.NodeElement {
				ui.data.Set(varVisible(n.(dom.ElementNode).ID()), false)
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
	renderNode(ctx, ui.document, ui.data, ui.jsvm)
}

func renderNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) {
	attrs := node.Attributes()
	sn, cn := pushStyles(attrs, data)
	defer popStyles(sn, cn)
	switch node.TagName() {
	case "_root":
		renderRootNode(ctx, node, data, jsvm)
	case "window":
		renderWindowNode(ctx, node, data, jsvm)
	case "demowindow":
		renderDemoWindowNode(ctx, node, data, jsvm)
	case "button":
		if imgui.Button(z.S(attrs["label"], node.FirstChildAsText())) {
			bubbles := true
			if jsv := attrs["onclick"]; jsv != "" {
				rv, err := jsvm.RunString(jsv)
				if err != nil {
					println("js err: " + err.Error())
					bubbles = false
				} else {
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
}

func renderRootNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) {
	for _, child := range node.Children() {
		if child.Type() == dom.NodeElement {
			renderNode(ctx, child.(dom.ElementNode), data, jsvm)
		} else if child.Type() == dom.NodeText {
			imgui.Text(strings.TrimSpace(child.(dom.TextNode).Text()))
		}
	}
}

func renderWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) {
	attrs := node.Attributes()
	if node.ID() == "" {
		println("warning: window didn't have an ID")
		node.SetAttribute("id", z.Rs())
	}
	wname := attrs["name"]
	if wname == "" {
		wname = node.ID()
	}
	show := data.MustBool(varVisible(node.ID()), node.Attributes().BoolD("visible", true))
	lshow := show
	if show {
		setupWindowPos(ctx, attrs, data, jsvm)
		imgui.BeginV(wname, &show, parseWindowFlags(attrs))
		for _, child := range node.Children() {
			if child.Type() == dom.NodeElement {
				renderNode(ctx, child.(dom.ElementNode), data, jsvm)
			} else if child.Type() == dom.NodeText {
				imgui.Text(child.(dom.TextNode).Text())
			}
		}
		imgui.End()
	}
	if lshow != show {
		// the window was closed by the X button
		data.Set(varVisible(node.ID()), show)
	}
}

func renderDemoWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) {
	if node.ID() == "" {
		println("warning: window didn't have an ID")
		node.SetAttribute("id", z.Rs())
	}
	show := data.MustBool(varVisible(node.ID()), node.Attributes().BoolD("visible", true))
	lshow := show
	if show {
		imgui.ShowDemoWindow(&show)
	}
	if lshow != show {
		// the window was closed by the X button
		data.Set(varVisible(node.ID()), show)
	}
}
