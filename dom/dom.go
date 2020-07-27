// Package dom provides methods and structs to construct and manipulate a simple
// Document Object Model
package dom

import (
	"strconv"
	"strings"
)

// NodeType represents an Element/Object or Text
type NodeType int

// Valid DOM nodes
const (
	NodeElement NodeType = iota
	NodeText    NodeType = iota
)

// Attributes of an Element node
type Attributes map[string]string

// BoolD retrieves the attribute [name]. Returns [defaultv] if not found
func (a Attributes) BoolD(name string, defaultv bool) bool {
	vs := a[name]
	if vs == "" {
		return defaultv
	}
	v, _ := strconv.ParseBool(vs)
	return v
}

// IntD retrieves the attribute [name]. Returns [defaultv] if not found
func (a Attributes) IntD(name string, defaultv int) int {
	vs := a[name]
	if vs == "" {
		return defaultv
	}
	v, _ := strconv.Atoi(vs)
	return v
}

func (a Attributes) String(name string) string {
	return a[name]
}

// HasAttr returns true if one of the attributes is found
func (a Attributes) HasAttr(name string, names ...string) bool {
	if _, ok := a[name]; ok {
		return true
	}
	for _, nm := range names {
		if _, ok := a[nm]; ok {
			return true
		}
	}
	return false
}

// FirstAttr returns the value of the first attibute that is found (if the name matches)
func (a Attributes) FirstAttr(name string, names ...string) string {
	if v, ok := a[name]; ok {
		return v
	}
	for _, nm := range names {
		if v, ok := a[nm]; ok {
			return v
		}
	}
	return ""
}

// Node is a valid DOM node
type Node interface {
	Type() NodeType
}

// TextNode is a valid DOM node with a text only content
type TextNode interface {
	Node
	Text() string
}

type ElementNode interface {
	Node
	TagName() string
	Children() []Node
	Attributes() Attributes
	ID() string
	Classes() []string
	Append(n Node)
	SetAttribute(name, value string)
	DeleteAttribute(name string)
	FirstChildAsText() string
	FindChildByID(id string) Node
}

type textNode struct {
	data string
}

var _ TextNode = (*textNode)(nil)

func (n *textNode) Type() NodeType {
	return NodeText
}

func (n *textNode) Text() string {
	return n.data
}

type elementNode struct {
	tagname    string
	attributes Attributes
	children   []Node
}

var _ ElementNode = (*elementNode)(nil)

func (n *elementNode) Type() NodeType {
	return NodeElement
}

func (n *elementNode) TagName() string {
	return n.tagname
}

func (n *elementNode) Children() []Node {
	return n.children
}

func (n *elementNode) Attributes() Attributes {
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

func (n *elementNode) SetAttribute(name, value string) {
	if n.attributes == nil {
		n.attributes = make(map[string]string)
	}
	n.attributes[name] = value
}

func (n *elementNode) DeleteAttribute(name string) {
	if n.attributes == nil {
		return
	}
	delete(n.attributes, name)
}

func (n *elementNode) FirstChildAsText() string {
	if len(n.children) > 0 && n.children[0].Type() == NodeText {
		return n.children[0].(TextNode).Text()
	}
	return ""
}

func (n *elementNode) FindChildByID(id string) Node {
	if len(n.children) < 1 {
		return nil
	}
	for _, child := range n.children {
		if child.Type() == NodeElement {
			if child.(ElementNode).ID() == id {
				return child
			}
		}
	}
	for _, child := range n.children {
		if child.Type() == NodeElement {
			if x := child.(ElementNode).FindChildByID(id); x != nil {
				return x
			}
		}
	}
	return nil
}

func Text(str string) TextNode {
	return &textNode{
		data: str,
	}
}

func Element(tagname string, attributes map[string]string, children ...Node) ElementNode {
	n := &elementNode{
		tagname:    tagname,
		attributes: Attributes(attributes),
		children:   make([]Node, 0, len(children)),
	}
	if n.attributes == nil {
		n.attributes = make(Attributes)
	}
	for _, v := range children {
		n.children = append(n.children, v)
	}
	return n
}
