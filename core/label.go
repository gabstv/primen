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

//FIXME: label is broken

type Label struct {
	text     string
	face     font.Face
	area     image.Point
	dborder  image.Point
	color    color.Color
	filter   ebiten.Filter
	disabled bool // if true, the DrawableSystem will not draw this
	//
	x           float64 // logical X position
	y           float64 // logical Y position
	angle       float64 // radians
	scaleX      float64 // logical X scale (1 = 100%)
	scaleY      float64 // logical Y scale (1 = 100%)
	originX     float64 // X origin (0 = left; 0.5 = center; 1 = right)
	originY     float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	offsetX     float64 // Text rendering offset X (pixel unit)
	offsetY     float64 // Text rendering offset Y (pixel unit)
	faceOffsetX int     // Text rendering offset X (pixel unit)
	faceOffsetY int     // Text rendering offset Y (pixel unit)
	//

	base      *ebiten.Image
	textdirty bool
	realSize  image.Point
	//
	opt ebiten.DrawImageOptions
}

func NewLabel() Label {
	return Label{
		textdirty: true,
	}
}

func (l *Label) SetText(t string) *Label {
	if t == l.text {
		// no change
		return l
	}
	redraw := false
	if l.base == nil {
		redraw = true
	} else if l.area.Eq(image.ZP) {
		// calc if image needs to be bigger
		ff := l.validFontFace()
		p := text.MeasureString(l.text, ff)
		p.X += l.dborder.X
		p.Y += l.dborder.Y
		if ww, hh := l.base.Size(); ww < p.X || hh < p.Y {
			redraw = true
		}
	}
	l.text = t
	if redraw {
		l.setupInnerImage()
		l.textdirty = true
	}
	return l
}

func (l *Label) SetOrigin(ox, oy float64) *Label {
	l.originX = ox
	l.originY = oy
	return l
}

func (l *Label) SetColor(c color.Color) *Label {
	l.color = c
	l.textdirty = true
	return l
}

func (l *Label) SetFilter(f ebiten.Filter) {
	if f == l.filter {
		// no change
		return
	}
	l.filter = f
	l.setupInnerImage()
	l.textdirty = true
}

// ComputedSize returns the width and height of the printable area of the label
func (l *Label) ComputedSize() image.Point {
	return text.MeasureString(l.text, l.validFontFace())
}

// FontFaceHeight returns the font height
func (l *Label) FontFaceHeight() int {
	return l.validFontFace().Metrics().Height.Round()
}

// ResetTextOffset sets the text offset to default (0, l.FontFaceHeight())
func (l *Label) ResetTextOffset() {
	l.faceOffsetX = 0
	l.faceOffsetY = l.FontFaceHeight()
}

func (l *Label) validFontFace() font.Face {
	if l.face != nil {
		return l.face
	}
	return rx.DefaultFontFace()
}

func (l *Label) compute() {
	if l.disabled {
		// exit early because this is not going to be drawed
		return
	}
	if !l.textdirty {
		return
	}
	l.renderText()
}

func (l *Label) setupInnerImage() {
	if l.area.Eq(image.ZP) {
		// dynamic
		ff := l.validFontFace()
		p := text.MeasureString(l.text, ff)
		l.base, _ = ebiten.NewImage(p.X+l.dborder.X, p.Y+l.dborder.Y, l.filter)
	} else {
		l.base, _ = ebiten.NewImage(l.area.X, l.area.Y, l.filter)
	}
}

func (l *Label) renderText() {
	if !l.textdirty {
		return
	}
	ff := l.validFontFace()
	l.base.Fill(color.Transparent)
	text.Draw(l.base, l.text, ff, l.faceOffsetX, l.faceOffsetY, l.color)
	l.realSize = text.MeasureString(l.text, ff)
}

func (l *Label) Draw(ctx DrawCtx, d *Drawable) {
	if d == nil {
		return
	}
	if l.disabled {
		return
	}
	if l.base == nil {
		return
	}
	g := d.G(l.scaleX, l.scaleY, l.angle, l.x, l.y)
	o := &l.opt
	o.GeoM.Reset()
	o.GeoM.Translate(applyOrigin(float64(l.realSize.X), l.originX)+l.offsetX, applyOrigin(float64(l.realSize.Y), l.originY)+l.offsetY)
	o.GeoM.Concat(g)
	//TODO: reimplement colormode and composite mode
	ctx.Renderer().DrawImageRaw(l.base, o)

	if DebugDraw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+float64(l.realSize.X), y0
		x2, y2 := x1, y1+float64(l.realSize.Y)
		x3, y3 := x2-float64(l.realSize.X), y2
		screen := ctx.Renderer().Screen()
		debugLineM(screen, o.GeoM, x0, y0, x1, y1, debugBoundsColor)
		debugLineM(screen, o.GeoM, x1, y1, x2, y2, debugBoundsColor)
		debugLineM(screen, o.GeoM, x2, y2, x3, y3, debugBoundsColor)
		debugLineM(screen, o.GeoM, x3, y3, x0, y0, debugBoundsColor)
		debugLineM(screen, g, -4, 0, 4, 0, debugPivotColor)
		debugLineM(screen, g, 0, -4, 0, 4, debugPivotColor)
	}
}

//go:generate ecsgen -n Label -p core -o label_component.go --component-tpl --vars "UUID=1A74D1BE-BBF7-44F4-AC8B-18A00208EB76" --vars "BeforeRemove=c.beforeRemove(e)" --vars "OnAdd=c.onAdd(e)"

func (c *LabelComponent) beforeRemove(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = nil
	}
}

func (c *LabelComponent) onAdd(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = c.Data(e)
	} else {
		SetDrawableComponentData(c.world, e, Drawable{
			drawer: c.Data(e),
		})
	}
}

//go:generate ecsgen -n DrawableLabel -p core -o label_drawablesystem.go --system-tpl --vars "Priority=10"  --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "OnResize=s.onResize()" --vars "UUID=70EC2F13-4C71-4A3F-9F6D-FF11F5DE9384" --components "Drawable" --components "Label"

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
	d := GetDrawableComponentData(s.world, e)
	d.drawer = GetLabelComponentData(s.world, e)
}

func (s *DrawableLabelSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *DrawableLabelSystem) onResize() {
	for _, v := range s.V().Matches() {
		v.Drawable.drawer = v.Label
	}
}

// DrawPriority noop
func (s *DrawableLabelSystem) DrawPriority(ctx DrawCtx) {

}

// Draw noop (drawing is controlled by *Drawable)
func (s *DrawableLabelSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *DrawableLabelSystem) UpdatePriority(ctx UpdateCtx) {}

// Update computes labes if dirty
func (s *DrawableLabelSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Label.compute()
	}
}

// func (l *Label) Destroy() {
// 	l.base = nil
// 	l.SetDirty()
// 	l.transformMatrix = nil
// }
