package ggui

import (
	"image"

	"github.com/gabstv/primen/core"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/internal/z"
	"github.com/hajimehoshi/ebiten"
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
	world    ecs.BaseWorld
	id       string
	disabled bool
	document dom.ElementNode
}

func NewUIManager() UIManager {
	return UIManager{}
}

func (m *UIManager) Setup(root dom.ElementNode) {
	m.document = root
	m.build()
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

func (m *UIManager) buildrect(elem dom.ElementNode, pentity ecs.Entity) ecs.Entity {
	entity := m.world.NewEntity()
	SetUINodeComponentData(m.world, entity, UINode{
		id:          z.S(elem.ID(), z.Rs()),
		uimanagerid: m.id,
	})
	attrs := elem.Attributes()
	SetUIRectComponentData(m.world, entity, UIRect{
		filter:      z.Filter(attrs["filter"], ebiten.FilterDefault),
		bgColor:     z.Color(z.S(attrs["bgcolor"], attrs["background-color"]), z.White),
		stroke:      z.Int(z.S(attrs["stroke-size"], attrs["strokesz"]), 0),
		strokeColor: z.Color(z.S(attrs["strokec"], attrs["stroke-color"]), z.Black),
		size: image.Point{
			X: z.Int(z.S(attrs["width"], attrs["w"]), 0),
			Y: z.Int(z.S(attrs["height"], attrs["h"]), 0),
		},
	})
	core.SetTransformComponentData(m.world, entity, core.Transform{})
	trtr := core.GetTransformComponentData(m.world, entity)
	trtr.SetX(z.Float64(attrs["x"], 0)).SetY(z.Float64(attrs["y"], 0))
	trtr.SetParent(pentity)
	trtr.SetScale(z.Float64(z.S(attrs["scalex"], attrs["sx"]), 1), z.Float64(z.S(attrs["scaley"], attrs["sy"]), 1))
	trtr.SetAngle(z.Float64(z.S(attrs["rotation"], attrs["rot"], attrs["angle"]), 0))
	//TODO: add other components (?)
	return entity
}

//go:generate ecsgen -n UIManager -p core -o uimanager_component.go --component-tpl --vars "UUID=D81D8469-5C53-4436-9323-74635C5BF624" --vars "Setup=c.onCompSetup()" --vars "OnAdd=c.setupNewComp(e)"

func (c *UIManagerComponent) onCompSetup() {

}

func (c *UIManagerComponent) setupNewComp(e ecs.Entity) {
	d := c.Data(e)
	d.world = c.world
}
