package troupe

import (
	"context"
	"time"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

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
	fn2 := func(ctx ecs.Context) {
		scr := w.Get("screen").(*ebiten.Image)
		ctx2 := ctx.(Context)
		fn(ctx2, scr)
	}
	return w.World.NewSystem(priority, fn2, comps...)
}

func (w *World) NewComponent(input NewComponentInput) (*Component, error) {
	return w.World.NewComponent(input)
}

func NewWorld(e *Engine) *World {
	w := &World{
		World: ecs.NewWorldWithCtx(func(c0 context.Context, dt float64, sys *System, w Worlder) ecs.Context {
			return ctxt{
				c:          c0,
				dt:         dt,
				system:     sys,
				world:      w,
				engine:     e,
				fps:        ebiten.CurrentFPS(),
				drwskipped: ebiten.IsDrawingSkipped(),
			}
		}),
	}
	return w
}
