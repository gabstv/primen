package ggebiten

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Graphics struct {
	GG     *gg.Context
	linked LinkedImage
}

func (g *Graphics) Sync() {
	g.linked.UpdatePixels()
}

func (g *Graphics) Dispose() {
	g.linked.Dispose()
}

func (g *Graphics) Ebimage() *ebiten.Image {
	return g.linked.Ebimage()
}

func f(v int) float64 {
	return float64(v)
}

func (g *Graphics) DrawRect(x, y, w, h, strokewidth int, stroke, bg color.Color) {
	c := g.GG
	if strokewidth > 0 {
		c.SetLineWidth(f(strokewidth))
		c.SetLineCap(gg.LineCapSquare)
		c.SetLineJoin(gg.LineJoinBevel)
	}
	c.SetFillRule(gg.FillRuleWinding)
	c.SetColor(bg)
	c.DrawRectangle(f(x), f(y), f(w), f(h))
	if strokewidth > 0 {
		c.FillPreserve()
		c.SetColor(stroke)
		c.Stroke()
	} else {
		c.Fill()
	}
}

func NewGraphicsSoftLink(width, height int, filter ebiten.Filter) *Graphics {
	gr := &Graphics{
		linked: NewSoftLinkedImage(width, height, filter),
	}
	gr.GG = gg.NewContextForRGBA(gr.linked.Image())
	return gr
}
