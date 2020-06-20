package core

import (
	"image"

	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// TileSet is a drawable that efficiently draws tiles in a 2d array
type TileSet struct {
	x        float64 // logical X position
	y        float64 // logical Y position
	angle    float64 // radians
	scaleX   float64 // logical X scale (1 = 100%)
	scaleY   float64 // logical Y scale (1 = 100%)
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
}

func (t *TileSet) Draw(ctx DrawCtx, d *Drawable) {
	if d == nil {
		return
	}
	if t.disabled {
		return
	}
	g := d.G(t.scaleX, t.scaleY, t.angle, t.x, t.y)
	o := &t.opt
	o.GeoM.Reset()
	o.GeoM.Translate(applyOrigin(t.cellWidth*float64(t.cSize.X), t.originX)+t.offsetX, applyOrigin(t.cellHeight*float64(t.cSize.Y), t.originY)+t.offsetY)
	o.GeoM.Concat(g)
	tilem := GeoM()
	if t.cSize.X <= 0 {
		// invalid tile size
		return
	}
	renderer := ctx.Renderer()
	for i, p := range t.cells {
		y := i / t.cSize.X
		x := i % t.cSize.X
		tilem.Reset().Translate(float64(x)*t.cellWidth, float64(y)*t.cellHeight)
		tilem.Concat(o.GeoM)
		renderer.DrawImage(t.db[p], tilem)
	}
	//TODO: debug draw
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

//go:generate ecsgen -n DrawableTileSet -p core -o tileset_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "OnResize=s.onResize()" --vars "UUID=D02C8E02-E7C0-497B-B332-91543DF4FFFD" --components "Drawable" --components "TileSet"

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
	d.drawer = GetTileSetComponentData(s.world, e)
}

func (s *DrawableTileSetSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *DrawableTileSetSystem) onResize() {
	for _, v := range s.V().Matches() {
		v.Drawable.drawer = v.TileSet
	}
}

// DrawPriority noop
func (s *DrawableTileSetSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop (controlled by *Drawable)
func (s *DrawableTileSetSystem) Draw(ctx DrawCtx) {

}

// UpdatePriority noop
func (s *DrawableTileSetSystem) UpdatePriority(ctx UpdateCtx) {}

// Update noop
//TODO: remove logic
func (s *DrawableTileSetSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		if v.TileSet.db == nil {
			v.TileSet.isValid = false
			continue
		}
		if len(v.TileSet.cells) < v.TileSet.cSize.X*v.TileSet.cSize.Y {
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
// 		g = GeoM().Scale(t.scaleX, t.scaleY).Rotate(t.Angle).Translate(t.X, t.Y)
// 	}
// 	lg := GeoM().Translate(applyOrigin(t.cellWidth*float64(t.cSize.X), t.originX), applyOrigin(t.cellHeight*float64(t.cSize.Y), t.originY))
// 	lg.Translate(t.offsetX, t.offsetY)
// 	lg.Concat(*g.M())
// 	tilem := GeoM()
// 	if t.cSize.X <= 0 {
// 		// invalid tile size
// 		return
// 	}
// 	for i, p := range t.cells {
// 		y := i / t.cSize.X
// 		x := i % t.cSize.X
// 		tilem.Reset().Translate(float64(x)*t.cellWidth, float64(y)*t.cellHeight)
// 		tilem.Concat(*lg.M())
// 		renderer.DrawImage(t.db[p], tilem)
// 	}
// }

// // Destroy Drawable implementation
// func (t *TileSet) Destroy() {
// 	//TODO: implement
// }

// // IsDisabled Drawable implementation
// func (t *TileSet) IsDisabled() bool {
// 	return t.disabled
// }

// // Size Drawable implementation
// func (t *TileSet) Size() (w, h float64) {
// 	return t.cellWidth * float64(t.cSize.X), t.cellHeight * float64(t.cSize.Y)
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
// 	t.offsetX = x
// 	t.offsetY = y
// }
