package imgui

import (
	"strings"

	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/css"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
	"github.com/inkyblackness/imgui-go/v2"
)

type UI struct {
	id       UID
	data     *uiMemory
	document dom.ElementNode
	styles   []*css.Stylesheet
	// css.Style
}

func newUI(id UID, doc dom.ElementNode, styles ...*css.Stylesheet) *UI {
	ui := &UI{
		id:       id,
		document: doc,
		styles:   styles,
		data:     newUIMemory(),
	}
	ui.inlineStyles()
	return ui
}

func (ui *UI) inlineStyles() {

}

func (ui *UI) Render(ctx core.DrawCtx) {
	renderNode(ctx, ui.document, ui.data)
}

func renderNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory) {
	attrs := node.Attributes()
	sn, cn := pushStyles(attrs, data)
	defer popStyles(sn, cn)
	switch node.TagName() {
	case "_root":
		renderRootNode(ctx, node, data)
	case "window":
		renderWindowNode(ctx, node, data)
	case "button":
		if imgui.Button(z.S(attrs["label"], node.FirstChildAsText())) {
			//ctx.Engine().
			println("TODO: button click!!!")
		}
	}
}

func renderRootNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory) {
	for _, child := range node.Children() {
		if child.Type() == dom.NodeElement {
			renderNode(ctx, child.(dom.ElementNode), data)
		} else if child.Type() == dom.NodeText {
			imgui.Text(strings.TrimSpace(child.(dom.TextNode).Text()))
		}
	}
}
func renderWindowNode(ctx core.DrawCtx, node dom.ElementNode, data *uiMemory) {
	attrs := node.Attributes()
	if node.ID() == "" {
		println("warning: window didn't have an ID")
		node.SetAttribute("id", z.Rs())
	}
	wname := attrs["name"]
	if wname == "" {
		wname = node.ID()
	}
	if *data.UpsertBool(node.ID()+"_active", true) {
		imgui.BeginV(wname, data.UpsertBool(node.ID()+"_active", true), parseWindowFlags(attrs))
		for _, child := range node.Children() {
			if child.Type() == dom.NodeElement {
				renderNode(ctx, child.(dom.ElementNode), data)
			} else if child.Type() == dom.NodeText {
				imgui.Text(child.(dom.TextNode).Text())
			}
		}
		imgui.End()
	}
}
