package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
	DrawTarget(id DrawTargetID) DrawTarget
}
