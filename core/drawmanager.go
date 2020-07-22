package core

import (
	"github.com/hajimehoshi/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
}

type drawManager struct {
	__screen *ebiten.Image
	imopt    *ebiten.DrawImageOptions
	prevgeom ebiten.GeoM
	prevcm   ebiten.ColorM
	prevmode ebiten.CompositeMode
}

func newDrawManager(screen *ebiten.Image) DrawManager {
	return &drawManager{
		__screen: screen,
		imopt:    &ebiten.DrawImageOptions{}, //TODO: default filter
	}
}

func (m *drawManager) screen() *ebiten.Image {
	return m.__screen
}

func (m *drawManager) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask DrawMask) {
	_ = m.screen().DrawImage(image, opt)
}

func (m *drawManager) Screen() *ebiten.Image {
	return m.screen()
}
