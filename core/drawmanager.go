package core

import (
	"github.com/hajimehoshi/ebiten"
)

type DrawManager interface {
	DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, drawmask DrawMask)
	Screen() *ebiten.Image
}
