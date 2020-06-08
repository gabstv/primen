package core

import (
	"image"
	"image/color"

	"github.com/gabstv/primen/rx"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Label struct {
	Text         string
	Face         font.Face
	Area         image.Point
	Color        color.Color
	Filter       ebiten.Filter
	DrawDisabled bool // if true, the DrawableSystem will not draw this
	//
	X           float64 // logical X position
	Y           float64 // logical Y position
	Angle       float64 // radians
	ScaleX      float64 // logical X scale (1 = 100%)
	ScaleY      float64 // logical Y scale (1 = 100%)
	OriginX     float64 // X origin (0 = left; 0.5 = center; 1 = right)
	OriginY     float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	OffsetX     float64 // Text rendering offset X (pixel unit)
	OffsetY     float64 // Text rendering offset Y (pixel unit)
	FaceOffsetX int     // Text rendering offset X (pixel unit)
	FaceOffsetY int     // Text rendering offset Y (pixel unit)
	//
	Options *ebiten.DrawImageOptions
	//

	base       *ebiten.Image
	lastBounds image.Rectangle
	lastFilter ebiten.Filter
	notdirty   bool
	lastText   string
	realSize   image.Point
	//
	transformMatrix ebiten.GeoM
	customMatrix    bool
}

func (l *Label) ComputedSize() image.Point {
	return text.MeasureString(l.Text, l.validFontFace())
}

// FontFaceHeight returns the font height
func (l *Label) FontFaceHeight() int {
	return l.validFontFace().Metrics().Height.Round()
}

// ResetTextOffset sets the text offset to default (0, l.FontFaceHeight())
func (l *Label) ResetTextOffset() {
	l.FaceOffsetX = 0
	l.FaceOffsetY = l.FontFaceHeight()
}

func (l *Label) dirty() bool {
	return !l.notdirty
}

func (l *Label) textChanged() bool {
	if l.dirty() {
		return true
	}
	return l.lastText != l.Text
}

func (l *Label) SetDirty() {
	l.notdirty = false
}

func (l *Label) setNotDirty() {
	l.notdirty = true
}

func (l *Label) validFontFace() font.Face {
	if l.Face != nil {
		return l.Face
	}
	return rx.DefaultFontFace()
}

func (l *Label) compute() {
	if l.DrawDisabled {
		// exit early because this is not going to be drawed
		return
	}
	if !l.Area.Eq(image.ZP) {
		l.computeFixedArea()
	} else {
		l.computeDynamicArea()
	}
	if l.textChanged() {
		l.lastText = l.Text
		ff := l.validFontFace()
		if l.lastText == "" {
			l.base.Fill(color.Transparent)
		}
		text.Draw(l.base, l.Text, ff, l.FaceOffsetX, l.FaceOffsetY, l.Color)
		l.realSize = text.MeasureString(l.Text, ff)
	}
	l.setNotDirty()
}

func (l *Label) computeFixedArea() {
	if !l.dirty() {
		return
	}
	if l.base == nil || l.lastFilter != l.Filter || l.base.Bounds().Eq(l.lastBounds) {
		l.base, _ = ebiten.NewImage(l.Area.X, l.Area.Y, l.Filter)
	} else {
		_ = l.base.Fill(color.Transparent)
	}
	l.lastFilter = l.Filter
	l.lastBounds = l.base.Bounds()
}

func (l *Label) computeDynamicArea() {
	if !l.dirty() {
		return
	}
	ff := l.validFontFace()
	p := text.MeasureString(l.Text, ff)
	if l.base == nil || l.lastFilter != l.Filter || l.base.Bounds().Eq(l.lastBounds) {
		l.base, _ = ebiten.NewImage(p.X, p.Y, l.Filter)
	} else {
		_ = l.base.Fill(color.Transparent)
	}
	l.lastFilter = l.Filter
	l.lastBounds = l.base.Bounds()
}

// implements drawable

// Update does some computation before drawing
func (l *Label) Update(ctx Context) {
	l.compute()
}

// Draw is called by the Drawable systems
func (l *Label) Draw(screen *ebiten.Image, opt *ebiten.DrawImageOptions) {
	if l.DrawDisabled {
		return
	}
	prevGeo := opt.GeoM
	if l.customMatrix {
		opt.GeoM = l.transformMatrix
	} else {
		opt.GeoM.Scale(l.ScaleX, l.ScaleY)
		opt.GeoM.Rotate(l.Angle)
		opt.GeoM.Translate(l.X, l.Y)
	}
	xxg := &ebiten.GeoM{}
	xxg.Translate(applyOrigin(float64(l.realSize.X), l.OriginX), applyOrigin(float64(l.realSize.Y), l.OriginY))
	xxg.Translate(l.OffsetX, l.OffsetY)
	xxg.Concat(opt.GeoM)
	centerM := opt.GeoM
	opt.GeoM = *xxg

	// finally draw text
	screen.DrawImage(l.base, opt)
	if DebugDraw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+float64(l.realSize.X), y0
		x2, y2 := x1, y1+float64(l.realSize.Y)
		x3, y3 := x2-float64(l.realSize.X), y2
		debugLineM(screen, opt.GeoM, x0, y0, x1, y1, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x1, y1, x2, y2, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x2, y2, x3, y3, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x3, y3, x0, y0, debugBoundsColor)
		debugLineM(screen, centerM, -4, 0, 4, 0, debugPivotColor)
		debugLineM(screen, centerM, 0, -4, 0, 4, debugPivotColor)
	}
	opt.GeoM = prevGeo
}

func (l *Label) Destroy() {
	l.base = nil
	l.SetDirty()
	l.Options = nil
}

func (l *Label) DrawImageOptions() *ebiten.DrawImageOptions {
	return l.Options
}

func (l *Label) IsDisabled() bool {
	return l.DrawDisabled
}

// Size returns the real size of the label
func (l *Label) Size() (w, h float64) {
	return float64(l.realSize.X), float64(l.realSize.Y)
}

// SetTransformMatrix is used by TransformSystem to set a custom transform
func (l *Label) SetTransformMatrix(m ebiten.GeoM) {
	l.transformMatrix = m
	l.customMatrix = true
}

func (l *Label) ClearTransformMatrix() {
	l.customMatrix = false
}

func (l *Label) SetOffset(x, y float64) {
	l.OffsetX = x
	l.OffsetY = y
}
