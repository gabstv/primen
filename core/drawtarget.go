package core

import "github.com/hajimehoshi/ebiten"

// DrawMask is a flag to choose the draw target(s) of a drawable component
type DrawMask uint64

const (
	DrawMaskDefault DrawMask = 1
	DrawMaskNone    DrawMask = 0
	DrawMaskAll     DrawMask = 0xffffffffffffffff
)

type DrawTarget interface {
	ID() uint64
	DrawMask() DrawMask
	Image() *ebiten.Image
	Draw(screen *ebiten.Image)
}

type drawTarget struct {
	id        uint64
	layerMask uint64
	img       *ebiten.Image
}
