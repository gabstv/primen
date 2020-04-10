package tau

import (
	"context"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

// Component
//   Data
// System

func UpsertComponent(w *ecs.World, comp ecs.NewComponentInput) *ecs.Component {
	if c := w.Component(comp.Name); c != nil {
		return c
	}
	x, err := w.NewComponent(comp)
	if err != nil {
		panic(err)
	}
	return x
}

type SystemExecFn func(ctx Context)

type ComponentSystem interface {
	SystemName() string
	SystemPriority() int
	SystemExec() SystemExecFn
	SystemTags() []string
	Components(w ecs.Worlder) []*ecs.Component
}

type BaseComponentSystem struct {
}

func (cs *BaseComponentSystem) SystemPriority() int {
	return 0
}

func (cs *BaseComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		// noop
	}
}

func (cs *BaseComponentSystem) SystemTags() []string {
	return []string{"update"}
}

// func (cs *BaseComponentSystem) Components(w *ecs.World) []*ecs.Component {
// 	return w.NewComponent(ecs.NewComponentInput{})
// }

// // Run all systems
// func (w *World) Run(screen *ebiten.Image, delta float64) (taken time.Duration) {
// 	w.Set("screen", screen)
// 	return w.World.Run(delta)
// }

// // RunWithTag runs all systems with the specified tag
// func (w *World) RunWithTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
// 	w.Set("screen", screen)
// 	return w.World.RunWithTag(tag, delta)
// }

// // RunWithoutTag runs all systems without the specified tag
// func (w *World) RunWithoutTag(tag string, screen *ebiten.Image, delta float64) (taken time.Duration) {
// 	w.Set("screen", screen)
// 	return w.World.RunWithoutTag(tag, delta)
// }

// NewWorld creates a new world
func NewWorld(e *Engine) *ecs.World {
	if e == nil {
		panic("engine can't be nil")
	}
	w := ecs.NewWorldWithCtx(func(c0 context.Context, dt float64, sys *ecs.System, w ecs.WorldDicter) ecs.Context {
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
	})
	w.Set(DefaultImageOptions, &ebiten.DrawImageOptions{})
	return w
}

func init() {
	ecs.DefaultSystemExecWrapper = func(w *ecs.World, fn ecs.SystemExec) ecs.SystemExec {
		fn2 := func(ctx ecs.Context) {
			ctx2 := ctx.(Context)
			fn(ctx2)
		}
		return fn2
	}
}
