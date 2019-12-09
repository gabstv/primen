package troupe

import (
	"time"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

type Context = ecs.Context
type System = ecs.System
type Component = ecs.Component
type View = ecs.View
type Entity = ecs.Entity
type NewComponentInput = ecs.NewComponentInput
type Worlder = ecs.Worlder

type World struct {
	*ecs.World
}

func (w *World) Run(screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.Run(delta)
}

func (w *World) RunWithTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.RunWithTag(tag, delta)
}

func (w *World) RunWithoutTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.RunWithoutTag(tag, delta)
}

type SystemFn func(ctx Context, screen *ebiten.Image)

func (w *World) NewSystem(priority int, fn SystemFn, comps ...*Component) *System {
	fn2 := func(ctx Context) {
		scr := w.Get("screen").(*ebiten.Image)
		fn(ctx, scr)
	}
	return w.World.NewSystem(priority, fn2, comps...)
}

func (w *World) NewComponent(input NewComponentInput) (*Component, error) {
	return w.World.NewComponent(input)
}

func NewWorld() *World {
	w := &World{
		World: ecs.NewWorld(),
	}
	return w
}
