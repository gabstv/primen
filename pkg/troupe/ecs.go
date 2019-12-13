package troupe

import (
	"context"
	"time"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

// System -> ecs.System
type System = ecs.System

// Component -> ecs.Component
type Component = ecs.Component

// View -> ecs.View
type View = ecs.View

// Entity -> ecs.Entity
type Entity = ecs.Entity

// NewComponentInput -> ecs.NewComponentInput
type NewComponentInput = ecs.NewComponentInput

// Worlder -> ecs.Worlder
type Worlder = ecs.Worlder

// Dicter -> ecs.Dicter
type Dicter = ecs.Dicter

// WorldDicter -> ecs.WorldDicter
type WorldDicter = ecs.WorldDicter

type World struct {
	*ecs.World
}

// Run all systems
func (w *World) Run(screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.Run(delta)
}

// RunWithTag runs all systems with the specified tag
func (w *World) RunWithTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.RunWithTag(tag, delta)
}

// RunWithoutTag runs all systems without the specified tag
func (w *World) RunWithoutTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
	w.Set("screen", screen)
	return w.World.RunWithoutTag(tag, delta)
}

// SystemFn is the loop function of a system
type SystemFn func(ctx Context, screen *ebiten.Image)

// NewSystem creates a new system
func (w *World) NewSystem(priority int, fn SystemFn, comps ...*Component) *System {
	fn2 := func(ctx ecs.Context) {
		scr := w.Get("screen").(*ebiten.Image)
		ctx2 := ctx.(Context)
		fn(ctx2, scr)
	}
	return w.World.NewSystem(priority, fn2, comps...)
}

// NewComponent creates a new component
func (w *World) NewComponent(input NewComponentInput) (*Component, error) {
	return w.World.NewComponent(input)
}

// NewWorld creates a new world
func NewWorld(e *Engine) *World {
	w := &World{
		World: ecs.NewWorldWithCtx(func(c0 context.Context, dt float64, sys *System, w WorldDicter) ecs.Context {
			return ctxt{
				c:          c0,
				dt:         dt,
				system:     sys,
				world:      w,
				engine:     e,
				fps:        ebiten.CurrentFPS(),
				frame:      e.frame,
				drwskipped: ebiten.IsDrawingSkipped(),
				imopt:      w.Get(DefaultImageOptions).(*ebiten.DrawImageOptions),
			}
		}),
	}
	return w
}
