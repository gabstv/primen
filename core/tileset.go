package core

import (
	"image"

	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// TileSet is a drawable that efficiently draws tiles in a 2d array
type TileSet struct {
	X            float64 // logical X position
	Y            float64 // logical Y position
	Angle        float64 // radians
	ScaleX       float64 // logical X scale (1 = 100%)
	ScaleY       float64 // logical Y scale (1 = 100%)
	OriginX      float64 // X origin (0 = left; 0.5 = center; 1 = right)
	OriginY      float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	OffsetX      float64 // offset origin X (in pixels)
	OffsetY      float64 // offset origin Y (in pixels)
	DrawDisabled bool    // if true, the DrawableSystem will not draw this

	DB         []*ebiten.Image
	CellWidth  float64
	CellHeight float64
	CSize      image.Point
	Cells      []int

	isValid         bool
	transformMatrix GeoMatrix
}

func (t *TileSet) Draw(ctx DrawCtx, o *Drawable) {
	g := &o.concatm
	if !o.concatset {
		g.Reset()
		g.Scale(t.ScaleX, t.ScaleY)
		g.Rotate(t.Angle)
		g.Translate(t.X, t.Y)
	}
	lg := GeoM().Translate(applyOrigin(t.CellWidth*float64(t.CSize.X), t.OriginX), applyOrigin(t.CellHeight*float64(t.CSize.Y), t.OriginY))
	lg.Translate(t.OffsetX, t.OffsetY)
	lg.Concat(*g)
	tilem := GeoM()
	if t.CSize.X <= 0 {
		// invalid tile size
		return
	}
	for i, p := range t.Cells {
		y := i / t.CSize.X
		x := i % t.CSize.X
		tilem.Reset().Translate(float64(x)*t.CellWidth, float64(y)*t.CellHeight)
		tilem.Concat(*lg.M())
		ctx.Renderer().DrawImage(t.DB[p], tilem)
	}
}

//go:generate ecsgen -n TileSet -p core -o tileset_component.go --component-tpl --vars "UUID=775FFA75-9F2F-423A-A905-D48E4D562AE8" --vars "BeforeRemove=c.beforeRemove(e)" --vars "OnAdd=c.onAdd(e)"

func (c *TileSetComponent) beforeRemove(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = nil
	}
}

func (c *TileSetComponent) onAdd(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = c.Data(e)
	} else {
		SetDrawableComponentData(c.world, e, Drawable{
			drawer: c.Data(e),
		})
	}
}

//go:generate ecsgen -n DrawableTileSet -p core -o tileset_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "UUID=D02C8E02-E7C0-497B-B332-91543DF4FFFD" --vars "ViewRescan=v.onRescan(e,x)" --components "Drawable" --components "TileSet"

func (v *viewDrawableTileSetSystem) onRescan(e ecs.Entity, x VIDrawableTileSetSystem) {
	x.Drawable.drawer = x.TileSet
}

var matchDrawableTileSetSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetTileSetComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchDrawableTileSetSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetTileSetComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *DrawableTileSetSystem) onEntityAdded(e ecs.Entity) {
	d := GetDrawableComponentData(s.world, e)
	d.Opt = &ebiten.DrawImageOptions{}
	d.drawer = GetTileSetComponentData(s.world, e)
}

func (s *DrawableTileSetSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *DrawableTileSetSystem) DrawPriority(ctx DrawCtx) {

}

func (s *DrawableTileSetSystem) Draw(ctx DrawCtx) {

}

// if Debug is TRUE
func (s *DrawableTileSetSystem) DebugDraw(ctx DrawCtx) {
	// screen := ctx.Renderer().Screen()
	// for _, v := range s.V().Matches() {
	// 	x0, y0 := 0.0, 0.0
	// 	x1, y1 := x0+v.Sprite.imageWidth, y0
	// 	x2, y2 := x1, y1+v.Sprite.imageHeight
	// 	x3, y3 := x2-v.Sprite.imageWidth, y2
	// 	debugLineM(screen, v.Drawable.Opt.GeoM, x0, y0, x1, y1, debugBoundsColor)
	// 	debugLineM(screen, v.Drawable.Opt.GeoM, x1, y1, x2, y2, debugBoundsColor)
	// 	debugLineM(screen, v.Drawable.Opt.GeoM, x2, y2, x3, y3, debugBoundsColor)
	// 	debugLineM(screen, v.Drawable.Opt.GeoM, x3, y3, x0, y0, debugBoundsColor)
	// 	debugLineM(screen, v.Drawable.concatm, -4, 0, 4, 0, debugPivotColor)
	// 	debugLineM(screen, v.Drawable.concatm, 0, -4, 0, 4, debugPivotColor)
	// }
}

func (s *DrawableTileSetSystem) UpdatePriority(ctx UpdateCtx) {

}

func (s *DrawableTileSetSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		if v.TileSet.DB == nil {
			v.TileSet.isValid = false
			continue
		}
		if len(v.TileSet.Cells) < v.TileSet.CSize.X*v.TileSet.CSize.Y {
			v.TileSet.isValid = false
			continue
		}
		v.TileSet.isValid = true
	}
	for _, v := range s.V().Matches() {
		if !v.TileSet.isValid {
			continue
		}
	}
}

// // Draw Drawable implementation
// func (t *TileSet) Draw(renderer DrawManager) {
// 	g := t.transformMatrix
// 	if g == nil {
// 		g = GeoM().Scale(t.ScaleX, t.ScaleY).Rotate(t.Angle).Translate(t.X, t.Y)
// 	}
// 	lg := GeoM().Translate(applyOrigin(t.CellWidth*float64(t.CSize.X), t.OriginX), applyOrigin(t.CellHeight*float64(t.CSize.Y), t.OriginY))
// 	lg.Translate(t.OffsetX, t.OffsetY)
// 	lg.Concat(*g.M())
// 	tilem := GeoM()
// 	if t.CSize.X <= 0 {
// 		// invalid tile size
// 		return
// 	}
// 	for i, p := range t.Cells {
// 		y := i / t.CSize.X
// 		x := i % t.CSize.X
// 		tilem.Reset().Translate(float64(x)*t.CellWidth, float64(y)*t.CellHeight)
// 		tilem.Concat(*lg.M())
// 		renderer.DrawImage(t.DB[p], tilem)
// 	}
// }

// // Destroy Drawable implementation
// func (t *TileSet) Destroy() {
// 	//TODO: implement
// }

// // IsDisabled Drawable implementation
// func (t *TileSet) IsDisabled() bool {
// 	return t.DrawDisabled
// }

// // Size Drawable implementation
// func (t *TileSet) Size() (w, h float64) {
// 	return t.CellWidth * float64(t.CSize.X), t.CellHeight * float64(t.CSize.Y)
// }

// // SetTransformMatrix Drawable implementation
// func (t *TileSet) SetTransformMatrix(m GeoMatrix) {
// 	t.transformMatrix = m
// }

// // ClearTransformMatrix Drawable implementation
// func (t *TileSet) ClearTransformMatrix() {
// 	t.transformMatrix = nil
// }

// // SetOffset Drawable implementation
// func (t *TileSet) SetOffset(x, y float64) {
// 	t.OffsetX = x
// 	t.OffsetY = y
// }
