// Package dom provides methods and structs to construct and manipulate a simple
// Document Object Model
package dom

import (
	"strings"
)

type NodeType int

const (
	NodeElement NodeType = iota
	NodeText    NodeType = iota
)

type Node interface {
	Type() NodeType
}

type TextNode interface {
	Node
	Text() string
}

type ElementNode interface {
	Node
	TagName() string
	Children() []Node
	Attributes() map[string]string
	ID() string
	Classes() []string
	Append(n Node)
}

type textNode struct {
	data string
}

func (n *textNode) Type() NodeType {
	return NodeText
}

func (n *textNode) Text() string {
	return n.data
}

type elementNode struct {
	tagname    string
	attributes map[string]string
	children   []Node
}

func (n *elementNode) Type() NodeType {
	return NodeElement
}

func (n *elementNode) TagName() string {
	return n.tagname
}

func (n *elementNode) Children() []Node {
	return n.children
}

func (n *elementNode) Attributes() map[string]string {
	return n.attributes
}

func (n *elementNode) ID() string {
	return n.attributes["id"]
}

func (n *elementNode) Classes() []string {
	return strings.Split(n.attributes["class"], " ")
}

func (n *elementNode) Append(node Node) {
	n.children = append(n.children, node)
}

func Text(str string) TextNode {
	return &textNode{
		data: str,
	}
}

func Element(tagname string, attributes map[string]string, children ...Node) ElementNode {
	n := &elementNode{
		tagname:    tagname,
		attributes: attributes,
		children:   make([]Node, 0, len(children)),
	}
	if n.attributes == nil {
		n.attributes = make(map[string]string)
	}
	for _, v := range children {
		n.children = append(n.children, v)
	}
	return n
}
