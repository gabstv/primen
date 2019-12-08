package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten"
)

func TestECS(t *testing.T) {
	w := NewWorld()
	c0, _ := w.NewComponent(NewComponentInput{
		Name: "test",
	})
	w.NewSystem(0, func(ctx Context, screen *ebiten.Image) {
		// do nothing
	}, c0)
	w.Run(nil, 1.0)
}
