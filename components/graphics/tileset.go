package graphics

import (
	"image"

	"github.com/gabstv/ebiten"
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
)

// TileSet is a drawable that efficiently draws tiles in a 2d array
type TileSet struct {
	originX  float64 // X origin (0 = left; 0.5 = center; 1 = right)
	originY  float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	offsetX  float64 // offset origin X (in pixels)
	offsetY  float64 // offset origin Y (in pixels)
	disabled bool    // if true, the DrawableSystem will not draw this

	db         []*ebiten.Image
	cellWidth  float64
	cellHeight float64
	cSize      image.Point
	cells      []int

	isValid bool
	opt     ebiten.DrawImageOptions

	drawMask core.DrawMask
}

func (s *TileSet) DrawMask() core.DrawMask {
	return s.drawMask
}

func (s *TileSet) SetDrawMask(mask core.DrawMask) {
	s.drawMask = mask
}

func NewTileSet(db []*ebiten.Image, rows, cols int, cellwidthpx, cellheightpx float64, cells []int) TileSet {
	tset := TileSet{
		drawMask: core.DrawMaskDefault,
		db:       db,
		cSize: image.Point{
			X: cols,
			Y: rows,
		},
		cellWidth:  cellwidthpx,
		cellHeight: cellheightpx,
		cells:      cells,
	}
	tset.isValid = tset.testvalid()
	return tset
}

func (t TileSet) testvalid() bool {
	if len(t.db) < 1 {
		return false
	}
	if t.cSize.X <= 0 || t.cSize.Y <= 0 {
		return false
	}
	if len(t.cells) < 1 {
		return false
	}
	return true
}

func (s *TileSet) SetEnabled(enabled bool) *TileSet {
	s.disabled = !enabled
	return s
}

func (s *TileSet) Origin() (ox, oy float64) {
	return s.originX, s.originY
}

func (s *TileSet) SetOrigin(ox, oy float64) *TileSet {
	s.originX, s.originY = ox, oy
	return s
}

func (s *TileSet) SetOffset(x, y float64) *TileSet {
	s.offsetX, s.offsetY = x, y
	return s
}

func (s *TileSet) SetDB(db []*ebiten.Image) *TileSet {
	s.db = db
	s.isValid = s.testvalid()
	return s
}

func (s *TileSet) SetColsRows(cols, rows int) *TileSet {
	s.cSize = image.Point{
		X: cols,
		Y: rows,
	}
	s.isValid = s.testvalid()
	return s
}

func (s *TileSet) SetCells(cells []int) *TileSet {
	s.cells = cells
	s.isValid = s.testvalid()
	return s
}

func (t *TileSet) Update(ctx core.UpdateCtx, tr *components.Transform) {
	if t.db == nil {
		t.isValid = false
		return
	}
	if len(t.cells) < t.cSize.X*t.cSize.Y {
		t.isValid = false
		return
	}
	t.isValid = true
}

func (t *TileSet) Draw(ctx core.DrawCtx, tr *components.Transform) {
	if t.disabled {
		return
	}
	g := tr.GeoM()
	o := &t.opt
	o.GeoM.Reset()
	o.GeoM.Translate(core.ApplyOrigin(t.cellWidth*float64(t.cSize.X), t.originX)+t.offsetX, core.ApplyOrigin(t.cellHeight*float64(t.cSize.Y), t.originY)+t.offsetY)
	o.GeoM.Concat(g)
	imopt := &ebiten.DrawImageOptions{}
	if t.cSize.X <= 0 {
		// invalid tile size
		return
	}
	renderer := ctx.Renderer()
	for i, p := range t.cells {
		y := i / t.cSize.X
		x := i % t.cSize.X
		imopt.GeoM.Reset()
		imopt.GeoM.Translate(float64(x)*t.cellWidth, float64(y)*t.cellHeight)
		imopt.GeoM.Concat(o.GeoM)
		renderer.DrawImage(t.db[p], imopt, t.drawMask)
	}
	//TODO: debug draw
}

//go:generate ecsgen -n TileSet -p graphics -o tileset_component.go --component-tpl --vars "UUID=775FFA75-9F2F-423A-A905-D48E4D562AE8" --vars "Setup=c.onCompSetup()"

func (c *TileSetComponent) onCompSetup() {
	RegisterDrawableComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) Drawable {
		return GetTileSetComponentData(w, e)
	})
}
