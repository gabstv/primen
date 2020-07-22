package core

import (
	"github.com/hajimehoshi/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
}

type drawManager struct {
	screen *ebiten.Image
	imopt  *ebiten.DrawImageOptions
}

func newDrawManager(screen *ebiten.Image) DrawManager {
	return &drawManager{
		screen: screen,
		imopt:  &ebiten.DrawImageOptions{}, //TODO: default filter
	}
}

func (m *drawManager) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask DrawMask) {
	_ = m.screen.DrawImage(image, opt)
}

func (m *drawManager) Screen() *ebiten.Image {
	return m.screen
}
