package renderer

import (
	"github.com/gabstv/primen/dom"
	"github.com/inkyblackness/imgui-go"
)

func Children(ctx *Context, node dom.ElementNode) {
	for _, child := range node.Children() {
		if child.Type() == dom.NodeElement {
			Node(ctx, child.(dom.ElementNode))
		} else if child.Type() == dom.NodeText {
			imgui.Text(child.(dom.TextNode).Text())
		}
	}
}
