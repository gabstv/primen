package primen

// import (
// 	"image"

// 	"github.com/gabstv/primen/core"
// 	"github.com/hajimehoshi/ebiten"
// )

// type TileMap struct {
// 	*WorldItem
// 	*TransformItem
// 	*DrawLayerItem
// 	tilemap func() *core.TileSet
// }

// func (t *TileMap) SetDB(db []*ebiten.Image) {
// 	t.tilemap().DB = db
// }

// func (t *TileMap) SetTilesXY(m2d [][]int) {
// 	if len(m2d) < 1 {
// 		return
// 	}
// 	width := len(m2d)
// 	height := len(m2d[0])
// 	mm := make([]int, width*height)
// 	for x := 0; x < width; x++ {
// 		for y := 0; y < height; y++ {
// 			mm[y*width+x] = m2d[x][y]
// 		}
// 	}
// 	t.tilemap().Cells = mm
// 	t.tilemap().CSize = image.Point{
// 		X: width,
// 		Y: height,
// 	}
// }

// func (t *TileMap) SetTilesYX(m2d [][]int) {
// 	if len(m2d) < 1 {
// 		return
// 	}
// 	height := len(m2d)
// 	width := len(m2d[0])
// 	mm := make([]int, width*height)
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			mm[y*width+x] = m2d[y][x]
// 		}
// 	}
// 	t.tilemap().Cells = mm
// 	t.tilemap().CSize = image.Point{
// 		X: width,
// 		Y: height,
// 	}
// }

// func (t *TileMap) SetCellSize(w, h float64) {
// 	t.tilemap().CellWidth = w
// 	t.tilemap().CellHeight = h
// }

// func (t *TileMap) SetOrigin(ox, oy float64) {
// 	t.tilemap().OriginX = ox
// 	t.tilemap().OriginY = oy
// }

// func NewTileSet(parent WorldTransform, layer Layer) *TileMap {
// 	w := parent.World()
// 	e := w.NewEntity()
// 	spr := &TileMap{}
// 	spr.WorldItem = newWorldItem(e, w)
// 	spr.TransformItem = newTransformItem(e, parent)
// 	spr.DrawLayerItem = newDrawLayerItem(e, w)
// 	core.SetTileSetComponentData(parent.World(), e, core.TileSet{
// 		ScaleX: 1,
// 		ScaleY: 1,
// 	})
// 	spr.tilemap = func() *core.TileSet { return core.GetTileSetComponentData(spr.world, e) }
// 	return spr
// }
