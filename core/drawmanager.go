package core

import (
	"github.com/gabstv/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
	DrawTarget(id DrawTargetID) DrawTarget
}
