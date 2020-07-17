package ggui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Button struct {
	cache *ebiten.Image
	opt   ebiten.DrawImageOptions
	size  image.Point
}

// a real button should have components like:
// Transform
// UINode
// UIInteractiveNode
// Button
