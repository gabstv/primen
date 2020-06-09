package core

import (
	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, x, y float64)
	DrawImageG(image *ebiten.Image, m GeoMatrix)
	DrawImageGC(image *ebiten.Image, m GeoMatrix, c ebiten.ColorM)
	DrawImageRaw(image *ebiten.Image, opt *ebiten.DrawImageOptions)
	Screen() *ebiten.Image
}

type drawManager struct {
	w        *ecs.World
	__screen *ebiten.Image
	imopt    *ebiten.DrawImageOptions
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

func (m *drawManager) DrawImage(image *ebiten.Image, x, y float64) {
	om := m.imopt.GeoM
	m.imopt.GeoM.Translate(x, y)
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = om
}

func (m *drawManager) DrawImageG(image *ebiten.Image, g GeoMatrix) {
	om := m.imopt.GeoM
	m.imopt.GeoM = *g.M()
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = om
}

func (m *drawManager) DrawImageGC(image *ebiten.Image, g GeoMatrix, c ebiten.ColorM) {
	om := m.imopt.GeoM
	oc := m.imopt.ColorM
	m.imopt.GeoM = *g.M()
	m.imopt.ColorM = c
	_ = m.screen().DrawImage(image, m.imopt)
	m.imopt.GeoM = om
	m.imopt.ColorM = oc
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

func (mm *gw) Reset() GeoMatrix {
	mm.m.Reset()
	return mm
}
