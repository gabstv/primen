package core

import (
	"github.com/hajimehoshi/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
}

type drawManager struct {
	screen  *ebiten.Image
	imopt   *ebiten.DrawImageOptions
	targets []DrawTarget
}

func newDrawManager(screen *ebiten.Image, rt ...DrawTarget) DrawManager {
	return &drawManager{
		screen:  screen,
		imopt:   &ebiten.DrawImageOptions{}, //TODO: default filter
		targets: rt,
	}
}

func (m *drawManager) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask DrawMask) {
	//TODO: check if assembly is doing branchless with this
	if len(m.targets) > 0 {
		for _, rt := range m.targets {
			rt.DrawImage(image, opt, mask)
		}
	} else {
		_ = m.screen.DrawImage(image, opt)
	}
}

func (m *drawManager) Screen() *ebiten.Image {
	return m.screen
}
