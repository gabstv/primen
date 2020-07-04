package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
)

type UINode struct {
	id          string
	uimanagerid string
}

//go:generate ecsgen -n UINode -p core -o uinode_component.go --component-tpl --vars "UUID=5ACA09B9-C488-46D7-A62E-5061C3DF5E6E"

type UIInteractiveNode struct {
	tabindex int
	links    UILinks
	active   bool
}

type UILinks struct {
	TopID    string
	LeftID   string
	RightID  string
	BottomID string
}

type UIManager struct {
	world    World
	id       string
	disabled bool
	document dom.ElementNode
}

func (m *UIManager) build() {
	rdom := m.document
	m.buildelement(rdom, 0)
}

func (m *UIManager) buildelement(elem dom.ElementNode, pentity ecs.Entity) {
	switch elem.TagName() {
	case "window":
		//TODO: build a window rect that is linked to the engine window size
		for _, child := range elem.Children() {
			if child.Type() == dom.NodeElement {
				m.buildelement(child.(dom.ElementNode), pentity)
			}
		}
	case "rect":
		m.buildrect(elem, pentity)
	}
}

func (m *UIManager) buildrect(elem dom.ElementNode, pentity ecs.Entity) {
	entity := m.world.NewEntity()
	SetUINodeComponentData(m.world, entity, UINode{
		id:          z.S(elem.ID(), z.Rs()),
		uimanagerid: m.id,
	})
	//TODO: add other components
	// add components
	_ = entity
}

//go:generate ecsgen -n UIManager -p core -o uimanager_component.go --component-tpl --vars "UUID=D81D8469-5C53-4436-9323-74635C5BF624" --vars "Setup=c.onCompSetup()"

func (c *UIManagerComponent) onCompSetup() {

}
