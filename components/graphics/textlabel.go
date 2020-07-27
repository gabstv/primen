package graphics

import (
	"image"
	"image/color"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/rx"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type TextLabel struct {
	text     string
	face     font.Face
	area     image.Point
	dborder  image.Point
	color    color.Color
	filter   ebiten.Filter
	disabled bool // if true, the DrawableSystem will not draw this
	//
	originX        float64 // X origin (0 = left; 0.5 = center; 1 = right)
	originY        float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	offsetX        float64 // Text rendering offset X (pixel unit)
	offsetY        float64 // Text rendering offset Y (pixel unit)
	faceOffsetAuto bool    // if true, apply font face height to faceOffsetY
	faceOffsetX    int     // Text rendering offset X (pixel unit)
	faceOffsetY    int     // Text rendering offset Y (pixel unit)
	//

	base      *ebiten.Image
	textdirty bool
	realSize  image.Point
	//
	opt ebiten.DrawImageOptions

	drawMask core.DrawMask
}

func (l *TextLabel) DrawMask() core.DrawMask {
	return l.drawMask
}

func (l *TextLabel) SetDrawMask(mask core.DrawMask) {
	l.drawMask = mask
}

func NewTextLabel() TextLabel {
	return TextLabel{
		drawMask:       core.DrawMaskDefault,
		textdirty:      true,
		faceOffsetAuto: true,
		color:          color.White,
	}
}

func (l *TextLabel) SetFaceOffset(x, y int) *TextLabel {
	l.faceOffsetX, l.faceOffsetY = x, y
	return l
}

func (l *TextLabel) SetFaceOffsetModeAuto(auto bool) *TextLabel {
	l.faceOffsetAuto = auto
	return l
}

func (l *TextLabel) SetText(t string) *TextLabel {
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
		p := text.MeasureString(t, ff)
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

func (l *TextLabel) SetOrigin(ox, oy float64) *TextLabel {
	l.originX = ox
	l.originY = oy
	return l
}

func (l *TextLabel) SetColor(c color.Color) *TextLabel {
	l.color = c
	l.textdirty = true
	return l
}

func (l *TextLabel) SetFilter(f ebiten.Filter) {
	if f == l.filter {
		// no change
		return
	}
	l.filter = f
	l.setupInnerImage()
	l.textdirty = true
}

// ComputedSize returns the width and height of the printable area of the label
func (l *TextLabel) ComputedSize() image.Point {
	return text.MeasureString(l.text, l.validFontFace())
}

// FontFaceHeight returns the font height
func (l *TextLabel) FontFaceHeight() int {
	return l.validFontFace().Metrics().Height.Round()
}

// ResetTextOffset sets the text offset to default (0, l.FontFaceHeight())
func (l *TextLabel) ResetTextOffset() {
	l.faceOffsetX = 0
	l.faceOffsetY = l.FontFaceHeight()
}

func (l *TextLabel) validFontFace() font.Face {
	if l.face != nil {
		return l.face
	}
	return rx.DefaultFontFace()
}

func (l *TextLabel) compute() {
	if l.disabled {
		// exit early because this is not going to be drawed
		return
	}
	if !l.textdirty {
		return
	}
	l.renderText()
}

func (l *TextLabel) setupInnerImage() {
	if l.area.Eq(image.ZP) {
		// dynamic
		ff := l.validFontFace()
		p := text.MeasureString(l.text, ff)
		l.base, _ = ebiten.NewImage(p.X+l.dborder.X, p.Y+l.dborder.Y, l.filter)
	} else {
		l.base, _ = ebiten.NewImage(l.area.X, l.area.Y, l.filter)
	}
}

func (l *TextLabel) renderText() {
	if !l.textdirty {
		return
	}
	ff := l.validFontFace()
	l.base.Fill(color.Transparent)
	autoh := 0
	if l.faceOffsetAuto {
		autoh = l.FontFaceHeight()
	}
	text.Draw(l.base, l.text, ff, l.faceOffsetX, l.faceOffsetY+autoh, l.color)
	l.realSize = text.MeasureString(l.text, ff)
	l.textdirty = false
}

func (l *TextLabel) Update(ctx core.UpdateCtx, tr *components.Transform) {}

func (l *TextLabel) Draw(ctx core.DrawCtx, tr *components.Transform) {
	if l.disabled {
		return
	}
	if l.base == nil {
		return
	}
	g := tr.GeoM()
	o := &l.opt
	o.GeoM.Reset()
	o.GeoM.Translate(core.ApplyOrigin(float64(l.realSize.X), l.originX)+l.offsetX, core.ApplyOrigin(float64(l.realSize.Y), l.originY)+l.offsetY)
	o.GeoM.Concat(g)
	//TODO: reimplement colormode and composite mode
	ctx.Renderer().DrawImage(l.base, o, l.drawMask)

	if debug.Draw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+float64(l.realSize.X), y0
		x2, y2 := x1, y1+float64(l.realSize.Y)
		x3, y3 := x2-float64(l.realSize.X), y2
		screen := ctx.Renderer().Screen()
		debug.LineM(screen, o.GeoM, x0, y0, x1, y1, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x1, y1, x2, y2, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x2, y2, x3, y3, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x3, y3, x0, y0, debug.BoundsColor)
		debug.LineM(screen, g, -4, 0, 4, 0, debug.PivotColor)
		debug.LineM(screen, g, 0, -4, 0, 4, debug.PivotColor)
	}
}

//go:generate ecsgen -n TextLabel -p graphics -o textlabel_component.go --component-tpl --vars "UUID=1A74D1BE-BBF7-44F4-AC8B-18A00208EB76"  --vars "Setup=c.onCompSetup()"

func (c *TextLabelComponent) onCompSetup() {
	RegisterDrawableComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) Drawable {
		return GetTextLabelComponentData(w, e)
	})
}

//go:generate ecsgen -n DrawableTextLabel -p graphics -o textlabel_drawablesystem.go --system-tpl --vars "Priority=10" --vars "UUID=70EC2F13-4C71-4A3F-9F6D-FF11F5DE9384" --components "TextLabel" --components "Transform;*components.Transform;components.GetTransformComponentData(v.world, e)" --go-import "\"github.com/gabstv/primen/components\""

var matchDrawableTextLabelSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(components.GetTransformComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetTextLabelComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchDrawableTextLabelSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(components.GetTransformComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetTextLabelComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *DrawableTextLabelSystem) DrawPriority(ctx core.DrawCtx) {

}

// Draw noop (drawing is controlled by *Drawable)
func (s *DrawableTextLabelSystem) Draw(ctx core.DrawCtx) {}

// UpdatePriority noop
func (s *DrawableTextLabelSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update computes labes if dirty
func (s *DrawableTextLabelSystem) Update(ctx core.UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.TextLabel.compute()
	}
}
