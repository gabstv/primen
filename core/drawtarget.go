package core

import "github.com/hajimehoshi/ebiten"

// DrawMask is a flag to choose the draw target(s) of a drawable component
type DrawMask uint64

const (
	// DrawMaskNone sets all bits to 0
	DrawMaskNone DrawMask = 0
	// DrawMaskAll sets all bits to 1
	DrawMaskAll DrawMask = 0xffffffffffffffff
)

var (
	// DrawMaskDefault can be changed to make all new drawable components
	// start with a different value.
	//
	// The initial value of DrawMaskDefault is DrawMaskAll
	DrawMaskDefault DrawMask = DrawMaskAll
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
