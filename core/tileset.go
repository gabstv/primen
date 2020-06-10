package core

import (
	"image"

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

// ██████╗  ██████╗   █████╗  ██╗    ██╗  █████╗  ██████╗  ██╗      ███████╗
// ██╔══██╗ ██╔══██╗ ██╔══██╗ ██║    ██║ ██╔══██╗ ██╔══██╗ ██║      ██╔════╝
// ██║  ██║ ██████╔╝ ███████║ ██║ █╗ ██║ ███████║ ██████╔╝ ██║      █████╗
// ██║  ██║ ██╔══██╗ ██╔══██║ ██║███╗██║ ██╔══██║ ██╔══██╗ ██║      ██╔══╝
// ██████╔╝ ██║  ██║ ██║  ██║ ╚███╔███╔╝ ██║  ██║ ██████╔╝ ███████╗ ███████╗
// ╚═════╝  ╚═╝  ╚═╝ ╚═╝  ╚═╝  ╚══╝╚══╝  ╚═╝  ╚═╝ ╚═════╝  ╚══════╝ ╚══════╝

// Update Drawable implementation
func (t *TileSet) Update(ctx Context) {
	if t.DB == nil {
		t.isValid = false
	}
	if len(t.Cells) < t.CSize.X*t.CSize.Y {
		t.isValid = false
	}
	t.isValid = true
}

// Draw Drawable implementation
func (t *TileSet) Draw(renderer DrawManager) {
	g := t.transformMatrix
	if g == nil {
		g = GeoM().Scale(t.ScaleX, t.ScaleY).Rotate(t.Angle).Translate(t.X, t.Y)
	}
	lg := GeoM().Translate(applyOrigin(t.CellWidth*float64(t.CSize.X), t.OriginX), applyOrigin(t.CellHeight*float64(t.CSize.Y), t.OriginY))
	lg.Translate(t.OffsetX, t.OffsetY)
	lg.Concat(*g.M())
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
		renderer.DrawImage(t.DB[p], tilem)
	}
}

// Destroy Drawable implementation
func (t *TileSet) Destroy() {
	//TODO: implement
}

// IsDisabled Drawable implementation
func (t *TileSet) IsDisabled() bool {
	return t.DrawDisabled
}

// Size Drawable implementation
func (t *TileSet) Size() (w, h float64) {
	return t.CellWidth * float64(t.CSize.X), t.CellHeight * float64(t.CSize.Y)
}

// SetTransformMatrix Drawable implementation
func (t *TileSet) SetTransformMatrix(m GeoMatrix) {
	t.transformMatrix = m
}

// ClearTransformMatrix Drawable implementation
func (t *TileSet) ClearTransformMatrix() {
	t.transformMatrix = nil
}

// SetOffset Drawable implementation
func (t *TileSet) SetOffset(x, y float64) {
	t.OffsetX = x
	t.OffsetY = y
}
