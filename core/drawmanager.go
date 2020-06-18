package core

import (
	"image/color"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

//FIXME: review

type DrawManager interface {
	DrawImageXY(image *ebiten.Image, x, y float64)
	DrawImage(image *ebiten.Image, m GeoMatrix)
	DrawImageC(image *ebiten.Image, m GeoMatrix, c ColorMatrix)
	DrawImageComp(image *ebiten.Image, m GeoMatrix, mode ebiten.CompositeMode)
	DrawImageCComp(image *ebiten.Image, m GeoMatrix, c ColorMatrix, mode ebiten.CompositeMode)
	DrawImageRaw(image *ebiten.Image, opt *ebiten.DrawImageOptions)
	Screen() *ebiten.Image
}

type drawManager struct {
	w        *ecs.World
	__screen *ebiten.Image
	imopt    *ebiten.DrawImageOptions
	prevgeom ebiten.GeoM
	prevcm   ebiten.ColorM
	prevmode ebiten.CompositeMode
}

func newDrawManager(w *ecs.World) DrawManager {
	return &drawManager{
		w:     w,
		imopt: &ebiten.DrawImageOptions{}, //TODO: default filter
	}
}

func (m *drawManager) screen() *ebiten.Image {
	if m.__screen != nil {
		return m.__screen
	}
	m.__screen = m.w.Get("screen").(*ebiten.Image)
	return m.__screen
}

func (m *drawManager) DrawImageXY(image *ebiten.Image, x, y float64) {
	om := m.imopt.GeoM
	m.imopt.GeoM.Translate(x, y)
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = om
}

func (m *drawManager) DrawImage(image *ebiten.Image, g GeoMatrix) {
	m.prevgeom = m.imopt.GeoM
	m.imopt.GeoM = g.MV()
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = m.prevgeom
}

func (m *drawManager) DrawImageC(image *ebiten.Image, g GeoMatrix, c ColorMatrix) {
	m.prevgeom = m.imopt.GeoM
	m.prevcm = m.imopt.ColorM
	m.imopt.GeoM = g.MV()
	m.imopt.ColorM = c.MV()
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = m.prevgeom
	m.imopt.ColorM = m.prevcm
}

func (m *drawManager) DrawImageComp(image *ebiten.Image, g GeoMatrix, mode ebiten.CompositeMode) {
	m.prevgeom = m.imopt.GeoM
	m.prevmode = m.imopt.CompositeMode
	m.imopt.GeoM = g.MV()
	m.imopt.CompositeMode = mode
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = m.prevgeom
	m.imopt.CompositeMode = m.prevmode
}

func (m *drawManager) DrawImageCComp(image *ebiten.Image, g GeoMatrix, c ColorMatrix, mode ebiten.CompositeMode) {
	m.prevgeom = m.imopt.GeoM
	m.prevcm = m.imopt.ColorM
	m.prevmode = m.imopt.CompositeMode
	m.imopt.GeoM = g.MV()
	m.imopt.ColorM = c.MV()
	m.imopt.CompositeMode = mode
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = m.prevgeom
	m.imopt.ColorM = m.prevcm
	m.imopt.CompositeMode = m.prevmode
}

func (m *drawManager) DrawImageRaw(image *ebiten.Image, opt *ebiten.DrawImageOptions) {
	_ = m.screen().DrawImage(image, opt)
}

func (m *drawManager) Screen() *ebiten.Image {
	return m.screen()
}

type gw struct {
	m *ebiten.GeoM
}

type GeoMatrix interface {
	Translate(tx, ty float64) GeoMatrix
	Scale(sx, sy float64) GeoMatrix
	Rotate(theta float64) GeoMatrix
	Concat(m ebiten.GeoM) GeoMatrix
	M() *ebiten.GeoM
	MV() ebiten.GeoM
	Reset() GeoMatrix
}

func GeoM() GeoMatrix {
	return &gw{
		m: &ebiten.GeoM{},
	}
}

func GeoM2(m ebiten.GeoM) GeoMatrix {
	return &gw{
		m: &m,
	}
}

func (mm *gw) Scale(sx, sy float64) GeoMatrix {
	mm.m.Scale(sx, sy)
	return mm
}

func (mm *gw) Rotate(theta float64) GeoMatrix {
	mm.m.Rotate(theta)
	return mm
}

func (mm *gw) Translate(tx, ty float64) GeoMatrix {
	mm.m.Translate(tx, ty)
	return mm
}

func (mm *gw) Concat(m ebiten.GeoM) GeoMatrix {
	mm.m.Concat(m)
	return mm
}

func (mm *gw) M() *ebiten.GeoM {
	return mm.m
}

func (mm *gw) MV() ebiten.GeoM {
	return *mm.m
}

func (mm *gw) Reset() GeoMatrix {
	mm.m.Reset()
	return mm
}

type ColorMatrix interface {
	Reset() ColorMatrix
	M() *ebiten.ColorM
	MV() ebiten.ColorM
	RotateHue(theta float64) ColorMatrix
	ChangeHSV(hueTheta float64, saturationScale float64, valueScale float64) ColorMatrix
	Apply(clr color.Color) ColorMatrix
}

type colorM struct {
	m *ebiten.ColorM
}

func (mm *colorM) Reset() ColorMatrix {
	mm.m.Reset()
	return mm
}

func (mm *colorM) RotateHue(theta float64) ColorMatrix {
	mm.m.RotateHue(theta)
	return mm
}

// ChangeHSV changes HSV (Hue-Saturation-Value) values.
// hueTheta is a radian value to rotate hue.
// saturationScale is a value to scale saturation.
// valueScale is a value to scale value (a.k.a. brightness).
//
// This conversion uses RGB to/from YCrCb conversion.
func (mm *colorM) ChangeHSV(hueTheta float64, saturationScale float64, valueScale float64) ColorMatrix {
	mm.m.ChangeHSV(hueTheta, saturationScale, valueScale)
	return mm
}

func (mm *colorM) Apply(clr color.Color) ColorMatrix {
	mm.m.Apply(clr)
	return mm
}

func (mm *colorM) M() *ebiten.ColorM {
	return mm.m
}

func (mm *colorM) MV() ebiten.ColorM {
	return *mm.m
}

func ColorM() ColorMatrix {
	return &colorM{
		m: &ebiten.ColorM{},
	}
}

func ColorM2(m ebiten.ColorM) ColorMatrix {
	return &colorM{
		m: &m,
	}
}

func ColorTint(c color.Color) ColorMatrix {
	mm := &colorM{
		m: &ebiten.ColorM{},
	}
	return mm.Apply(c)
}
