package primen

import (
	"github.com/gabstv/primen/components/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type TileSetNode struct {
	*Node
	wdl graphics.WatchDrawLayer
	wts graphics.WatchTileSet
}

func NewRootTileSetNode(w World, layer Layer, db []*ebiten.Image, rows, cols int, cellwidthpx, cellheightpx float64, cells []int) *TileSetNode {
	tsn := &TileSetNode{
		Node: NewRootNode(w),
	}
	graphics.SetDrawLayerComponentData(w, tsn.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetTileSetComponentData(w, tsn.e, graphics.NewTileSet(db, rows, cols, cellwidthpx, cellheightpx, cells))
	tsn.wdl = graphics.WatchDrawLayerComponentData(w, tsn.e)
	tsn.wts = graphics.WatchTileSetComponentData(w, tsn.e)
	return tsn
}

func NewChildTileSetNode(parent ObjectContainer, layer Layer, db []*ebiten.Image, rows, cols int, cellwidthpx, cellheightpx float64, cells []int) *TileSetNode {
	w := parent.World()
	tsn := &TileSetNode{
		Node: NewChildNode(parent),
	}
	graphics.SetDrawLayerComponentData(w, tsn.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetTileSetComponentData(w, tsn.e, graphics.NewTileSet(db, rows, cols, cellwidthpx, cellheightpx, cells))
	tsn.wdl = graphics.WatchDrawLayerComponentData(w, tsn.e)
	tsn.wts = graphics.WatchTileSetComponentData(w, tsn.e)
	return tsn
}

func (n *TileSetNode) TileSet() *graphics.TileSet {
	return n.wts.Data()
}

func (n *TileSetNode) SetLayer(l Layer) {
	n.wdl.Data().Layer = l
}

func (n *TileSetNode) SetZIndex(index int64) {
	n.wdl.Data().ZIndex = index
}

func (n *TileSetNode) Layer() Layer {
	return n.wdl.Data().Layer
}

func (n *TileSetNode) ZIndex() int64 {
	return n.wdl.Data().ZIndex
}
