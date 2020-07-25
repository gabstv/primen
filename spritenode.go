package primen

import (
	"image/color"

	"github.com/gabstv/ebiten"
	"github.com/gabstv/primen/components/graphics"
)

var (
	transparentPixel *ebiten.Image
)

type SpriteNode struct {
	*Node
	wdl  graphics.WatchDrawLayer
	wspr graphics.WatchSprite
}

func NewRootSpriteNode(w World, layer Layer) *SpriteNode {
	sprn := &SpriteNode{
		Node: NewRootNode(w),
	}
	graphics.SetDrawLayerComponentData(w, sprn.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetSpriteComponentData(w, sprn.e, graphics.NewSprite(0, 0, transparentPixel))
	sprn.wdl = graphics.WatchDrawLayerComponentData(w, sprn.e)
	sprn.wspr = graphics.WatchSpriteComponentData(w, sprn.e)
	return sprn
}

func NewChildSpriteNode(parent ObjectContainer, layer Layer) *SpriteNode {
	sprn := &SpriteNode{
		Node: NewChildNode(parent),
	}
	graphics.SetDrawLayerComponentData(parent.World(), sprn.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetSpriteComponentData(parent.World(), sprn.e, graphics.NewSprite(0, 0, transparentPixel))
	sprn.wdl = graphics.WatchDrawLayerComponentData(parent.World(), sprn.e)
	sprn.wspr = graphics.WatchSpriteComponentData(parent.World(), sprn.e)
	return sprn
}

func (n *SpriteNode) Sprite() *graphics.Sprite {
	return n.wspr.Data()
}

func (n *SpriteNode) SetLayer(l Layer) {
	n.wdl.Data().Layer = l
}

func (n *SpriteNode) SetZIndex(index int64) {
	n.wdl.Data().ZIndex = index
}

func (n *SpriteNode) Layer() Layer {
	return n.wdl.Data().Layer
}

func (n *SpriteNode) ZIndex() int64 {
	return n.wdl.Data().ZIndex
}

func init() {
	transparentPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = transparentPixel.Fill(color.Transparent)
}
