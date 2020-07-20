package core

import "github.com/hajimehoshi/ebiten"

type RenderTarget interface {
	ID() uint64
	LayerMask() uint64
	Image() *ebiten.Image
}

type renderTarget struct {
	id        uint64
	layerMask uint64
	img       *ebiten.Image
}
