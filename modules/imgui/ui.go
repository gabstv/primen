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

func (ui *UI) ID() UID {
	return ui.id
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
