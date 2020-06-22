package primen

import (
	"image/color"

	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
)

var (
	transparentPixel *ebiten.Image
)

type SpriteNode struct {
	*Node
	wdl  core.WatchDrawLayer
	wspr core.WatchSprite
}

func NewRootSpriteNode(w World, layer Layer) *SpriteNode {
	sprn := &SpriteNode{
		Node: NewRootNode(w),
	}
	core.SetDrawLayerComponentData(w, sprn.e, core.DrawLayer{
		Layer: layer,
	})
	core.SetSpriteComponentData(w, sprn.e, core.NewSprite(0, 0, transparentPixel))
	sprn.wdl = core.WatchDrawLayerComponentData(w, sprn.e)
	sprn.wspr = core.WatchSpriteComponentData(w, sprn.e)
	return sprn
}

func NewChildSpriteNode(parent ObjectContainer, layer Layer) *SpriteNode {
	sprn := &SpriteNode{
		Node: NewChildNode(parent),
	}
	core.SetDrawLayerComponentData(parent.World(), sprn.e, core.DrawLayer{
		Layer:  layer,
		ZIndex: core.ZIndexTop,
	})
	core.SetSpriteComponentData(parent.World(), sprn.e, core.NewSprite(0, 0, transparentPixel))
	sprn.wdl = core.WatchDrawLayerComponentData(parent.World(), sprn.e)
	sprn.wspr = core.WatchSpriteComponentData(parent.World(), sprn.e)
	return sprn
}

func (n *SpriteNode) Sprite() *core.Sprite {
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
