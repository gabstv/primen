package tau

import (
	"context"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

// Component
//   Data
// System

func UpsertComponent(w ecs.Worlder, comp ecs.NewComponentInput) *ecs.Component {
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
type SystemInitFn func(w *ecs.World, sys *ecs.System)

type ComponentSystem interface {
	SystemName() string
	SystemPriority() int
	SystemInit() SystemInitFn
	SystemExec() SystemExecFn
	SystemTags() []string
	Components(w ecs.Worlder) []*ecs.Component
}

type BaseComponentSystem struct {
}

func (cs *BaseComponentSystem) SystemPriority() int {
	return 0
}

func (cs *BaseComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		// noop
	}
}

func (cs *BaseComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		// noop
	}
}

func (cs *BaseComponentSystem) SystemTags() []string {
	return []string{"update"}
}

func SetupSystem(w *ecs.World, cs ComponentSystem) {
	fnfn := cs.SystemExec()
	wexec := func(ctx ecs.Context) {
		fnfn(ctx.(Context))
	}
	sys := w.NewSystem(cs.SystemName(), cs.SystemPriority(), wexec, cs.Components(w)...)
	if xinit := cs.SystemInit(); xinit != nil {
		xinit(w, sys)
	}
}

type BasicCS struct {
	SysName       string
	SysPriority   int
	SysInit       SystemInitFn
	SysExec       SystemExecFn
	SysTags       []string
	GetComponents func(w ecs.Worlder) []*ecs.Component
}

func (cs *BasicCS) SystemName() string {
	return cs.SysName
}

func (cs *BasicCS) SystemPriority() int {
	return cs.SysPriority
}

func (cs *BasicCS) BasicCS() int {
	return cs.SysPriority
}

func (cs *BasicCS) SystemInit() SystemInitFn {
	return cs.SysInit
}

func (cs *BasicCS) SystemExec() SystemExecFn {
	return cs.SysExec
}

func (cs *BasicCS) SystemTags() []string {
	return cs.SysTags
}

func (cs *BasicCS) Components(w ecs.Worlder) []*ecs.Component {
	return cs.GetComponents(w)
}

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
