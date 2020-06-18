package core

import (
	"image"
	"image/color"

	"github.com/gabstv/ecs/v2"
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

	base       *ebiten.Image
	lastBounds image.Rectangle
	lastFilter ebiten.Filter
	notdirty   bool
	lastText   string
	realSize   image.Point
	//
	transformMatrix GeoMatrix
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

//go:generate ecsgen -n Label -p core -o label_component.go --component-tpl --vars "UUID=1A74D1BE-BBF7-44F4-AC8B-18A00208EB76"

//go:generate ecsgen -n DrawableLabel -p core -o label_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "UUID=70EC2F13-4C71-4A3F-9F6D-FF11F5DE9384" --components "Drawable" --components "Label"

var matchDrawableLabelSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetLabelComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchDrawableLabelSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetLabelComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *DrawableLabelSystem) onEntityAdded(e ecs.Entity) {
	GetDrawableComponentData(s.world, e).opt = &ebiten.DrawImageOptions{}
}

func (s *DrawableLabelSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *DrawableLabelSystem) DrawPriority(ctx DrawCtx) {

}

func (s *DrawableLabelSystem) Draw(ctx DrawCtx) {

}

// if Debug is TRUE
func (s *DrawableLabelSystem) DebugDraw(ctx DrawCtx) {

}

func (s *DrawableLabelSystem) UpdatePriority(ctx UpdateCtx) {

}

func (s *DrawableLabelSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Label.compute()
		v.Drawable.img = v.Label.base
		v.Drawable.opt.GeoM.Reset()
		v.Drawable.opt.GeoM.Translate(applyOrigin(float64(v.Label.realSize.X), v.Label.OriginX)+v.Label.OffsetX, applyOrigin(float64(v.Label.realSize.Y), v.Label.OriginY)+v.Label.OffsetY)
		if v.Drawable.concatset {
			v.Drawable.opt.GeoM.Concat(v.Drawable.concatm)
		} else {
			v.Drawable.concatm.Reset()
			v.Drawable.concatm.Scale(v.Label.ScaleX, v.Label.ScaleY)
			v.Drawable.concatm.Rotate(v.Label.Angle)
			v.Drawable.concatm.Translate(v.Label.X, v.Label.Y)
			v.Drawable.opt.GeoM.Concat(v.Drawable.concatm)
		}
	}
}

// implements drawable

// // Update does some computation before drawing
// func (l *Label) Update(ctx Context) {
// 	l.compute()
// }

// // Draw is called by the Drawable systems
// func (l *Label) Draw(renderer DrawManager) {
// 	if l.DrawDisabled {
// 		return
// 	}
// 	g := l.transformMatrix
// 	if g == nil {
// 		g = GeoM().Scale(l.ScaleX, l.ScaleY).Rotate(l.Angle).Translate(l.X, l.Y)
// 	}
// 	lg := GeoM().Translate(applyOrigin(float64(l.realSize.X), l.OriginX), applyOrigin(float64(l.realSize.Y), l.OriginY))
// 	lg.Translate(l.OffsetX, l.OffsetY)
// 	lg.Concat(*g.M())
// 	renderer.DrawImage(l.base, lg)
// 	if DebugDraw {
// 		x0, y0 := 0.0, 0.0
// 		x1, y1 := x0+float64(l.realSize.X), y0
// 		x2, y2 := x1, y1+float64(l.realSize.Y)
// 		x3, y3 := x2-float64(l.realSize.X), y2
// 		debugLineM(renderer.Screen(), *lg.M(), x0, y0, x1, y1, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *lg.M(), x1, y1, x2, y2, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *lg.M(), x2, y2, x3, y3, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *lg.M(), x3, y3, x0, y0, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *g.M(), -4, 0, 4, 0, debugPivotColor)
// 		debugLineM(renderer.Screen(), *g.M(), 0, -4, 0, 4, debugPivotColor)
// 	}
// }

// func (l *Label) Destroy() {
// 	l.base = nil
// 	l.SetDirty()
// 	l.transformMatrix = nil
// }

// func (l *Label) IsDisabled() bool {
// 	return l.DrawDisabled
// }

// // Size returns the real size of the label
// func (l *Label) Size() (w, h float64) {
// 	return float64(l.realSize.X), float64(l.realSize.Y)
// }

// // SetTransformMatrix is used by TransformSystem to set a custom transform
// func (l *Label) SetTransformMatrix(m GeoMatrix) {
// 	l.transformMatrix = m
// }

// func (l *Label) ClearTransformMatrix() {
// 	l.transformMatrix = nil
// }

// func (l *Label) SetOffset(x, y float64) {
// 	l.OffsetX = x
// 	l.OffsetY = y
// }

// var _ Drawable = &Label{}
