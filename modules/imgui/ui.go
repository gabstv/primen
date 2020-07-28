package imgui

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gabstv/primen/core/js"

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

func newUI(id UID, doc []dom.Node) *UI {
	//goja.New()
	ui := &UI{
		id:       id,
		document: dom.Element("_root", nil, doc...),
		data:     newUIMemory(),
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
	lctx := setNodeLayout(node, data, jsvm)
	sn, cn := pushStyles(attrs, data)
	setNodeLayout(node, data, jsvm)
	defer popStyles(sn, cn)
	switch node.TagName() {
	case "_root":
		renderRootNode(ctx, node, data, jsvm, lctx)
	case "window":
		renderWindowNode(ctx, node, data, jsvm, lctx)
	case "demowindow":
		renderDemoWindowNode(ctx, node, data, jsvm, lctx)
	case "group":
		renderGroupNode(ctx, node, data, jsvm, lctx)
	case "columns":
		renderGroupColumns(ctx, node, data, jsvm, lctx)
	case "column":
		renderGroupColumn(ctx, node, data, jsvm, lctx)
	case "separator":
		imgui.Separator()
	case "spacing":
		imgui.Spacing()
	case "button":
		renderButtonNode(ctx, node, data, jsvm, lctx)
	}
}

func renderButtonNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
	attrs := node.Attributes()
	if imgui.ButtonV(z.S(attrs["label"], node.FirstChildAsText()), lctx.Size) {
		bubbles := true
		if jsv := attrs["onclick"]; jsv != "" {
			rv, err := jsvm.RunString(jsv)
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

func renderRootNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
	renderChildren(ctx, node, data, jsvm)
}

func renderWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
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
		renderChildren(ctx, node, data, jsvm)
		imgui.End()
	}
	if lshow != show {
		// the window was closed by the X button
		data.Set(varVisible(node.ID()), show)
	}
}

func renderGroupNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
	imgui.BeginGroup()
	renderChildren(ctx, node, data, jsvm)
	imgui.EndGroup()
}

func renderGroupColumns(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
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
	renderChildren(ctx, node, data, jsvm)
}

func renderGroupColumn(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
	renderChildren(ctx, node, data, jsvm)
	if lctx.Size.X > 0 {
		imgui.SetColumnWidth(-1, lctx.Size.X)
	}
	imgui.NextColumn()
}

func renderChildren(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) {
	for _, child := range node.Children() {
		if child.Type() == dom.NodeElement {
			renderNode(ctx, child.(dom.ElementNode), data, jsvm)
		} else if child.Type() == dom.NodeText {
			imgui.Text(child.(dom.TextNode).Text())
		}
	}
}

func renderDemoWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime, lctx layoutContext) {
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

func vhasJS(v string) (bool, string) {
	if strings.HasPrefix(v, "js:") {
		return true, v[3:]
	}
	return false, ""
}

func updateLocalJS(inited *bool, jsvm *goja.Runtime) {
	if *inited {
		return
	}
	jsvm.Set("__column_width", imgui.ColumnWidth())
	*inited = true
}

func getattrparsejsInt(jsvm *goja.Runtime, inited *bool, attrv string) (val int, relative, ok bool) {
	if attrv == "" {
		return 0, false, false
	}
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return 0, false, false
		}
		return int(vv.ToInteger()), false, true
	}
	if strings.HasSuffix(attrv, "%") {
		relative = true
		attrv = attrv[:len(attrv)-1]
	}
	vv, err := strconv.Atoi(attrv)
	if err != nil {
		return 0, relative, false
	}
	return vv, relative, true
}

func getattrparsejsNumber(jsvm *goja.Runtime, inited *bool, attrv string) (val float64, relative, ok bool) {
	if attrv == "" {
		return 0, false, false
	}
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return 0, false, false
		}
		return vv.ToFloat(), false, true
	}
	if strings.HasSuffix(attrv, "%") {
		relative = true
		attrv = attrv[:len(attrv)-1]
	}
	vv, err := strconv.ParseFloat(attrv, 64)
	if err != nil {
		return 0, relative, false
	}
	return vv, relative, true
}

func getattrparsejsString(jsvm *goja.Runtime, inited *bool, attrv string) (string, bool) {
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return "", false
		}
		return vv.String(), true
	}
	return attrv, true
}

type layoutContext struct {
	Size imgui.Vec2
}

func (c layoutContext) SizeZero() bool {
	return c.Size.X == 0 && c.Size.Y == 0
}

func setNodeLayout(node dom.ElementNode, data *uiMemory, jsvm *goja.Runtime) (layout layoutContext) {
	initjs := false
	attrs := node.Attributes()
	//if attrs.HasAttr("w", "width")
	if v, rel, ok := getattrparsejsNumber(jsvm, &initjs, attrs.FirstAttr("w", "width")); ok {
		if rel {
			cravail := imgui.ContentRegionAvail()
			layout.Size.X = cravail.X * (float32(v) / 100.0)
		} else {
			layout.Size.X = float32(v)
		}
	}
	if v, rel, ok := getattrparsejsNumber(jsvm, &initjs, attrs.FirstAttr("h", "height")); ok {
		if rel {
			cravail := imgui.ContentRegionAvail()
			layout.Size.X = cravail.Y * (float32(v) / 100.0)
		} else {
			layout.Size.X = float32(v)
		}
	}
	return
}
