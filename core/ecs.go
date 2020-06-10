package core

import (
	"context"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

var emptyCompSlice []*ecs.Component = make([]*ecs.Component, 0)

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
type SystemInitFn func(w *ecs.World, sys *ecs.System)

// Middleware is a system middleware
type Middleware func(next SystemExecFn) SystemExecFn

type ComponentSystem interface {
	SystemName() string
	SystemPriority() int
	SystemInit() SystemInitFn
	SystemExec() SystemExecFn
	SystemTags() []string
	Components(w *ecs.World) []*ecs.Component
	ExcludeComponents(w *ecs.World) []*ecs.Component
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

func (cs *BaseComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

func SetupSystem(w *ecs.World, cs ComponentSystem) {
	fnfn := cs.SystemExec()
	wexec := func(ctx ecs.Context) {
		fnfn(ctx.(Context))
	}
	sys := w.NewSystemX(cs.SystemName(), cs.SystemPriority(), wexec, cs.Components(w), cs.ExcludeComponents(w))
	for _, v := range cs.SystemTags() {
		sys.AddTag(v)
	}
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
	GetComponents func(w *ecs.World) []*ecs.Component
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

func (cs *BasicCS) Components(w *ecs.World) []*ecs.Component {
	return cs.GetComponents(w)
}

func (cs *BasicCS) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

// NewWorld creates a new world
func NewWorld(e Engine) *ecs.World {
	if e == nil {
		panic("engine can't be nil")
	}
	w := ecs.NewWorldWithCtx(func(c0 context.Context, dt float64, sys *ecs.System, w *ecs.World) ecs.Context {
		return ctxt{
			c:          c0,
			dt:         dt,
			system:     sys,
			world:      w,
			engine:     e,
			fps:        ebiten.CurrentFPS(),
			frame:      e.UpdateFrame(),
			drwskipped: ebiten.IsDrawingSkipped(),
			drawM:      newDrawManager(w),
		}
	})
	w.Set(DefaultImageOptions, &ebiten.DrawImageOptions{})
	return w
}

// SystemWrap wraps middlewares into a SystemExecFn
func SystemWrap(fn SystemExecFn, mid ...Middleware) SystemExecFn {
	for _, m := range mid {
		fn = m(fn)
	}
	return fn
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
