package core

import (
	"strconv"

	"github.com/gabstv/ebiten"
	"github.com/gabstv/primen/geom"
)

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

type DrawTargetID uint64

func (d DrawTargetID) String() string {
	return "DrawTarget#" + strconv.FormatUint(uint64(d), 10)
}

type DrawTarget interface {
	ID() DrawTargetID
	DrawMask() DrawMask
	Image() *ebiten.Image
	//DrawToScreen(screen *ebiten.Image)
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask DrawMask)
	Size() geom.Vec
	Translate(v geom.Vec)
	Scale(v geom.Vec)
	Rotate(rad float64)
	ResetTransform()
}
