package core

import (
	"math"

	"github.com/gabstv/ecs"
)

const (
	CNFunc       = "tau/core.FuncComponent"
	SNFunc       = "tau/core.FuncSystem"
	SNFuncAfter  = "tau/core.FuncAfterSystem"
	SNFuncBefore = "tau/core.FuncBeforeSystem"
)

type UpdateFn func(ctx Context, e ecs.Entity)

type Func struct {
	Fn       UpdateFn
	AfterFn  UpdateFn
	BeforeFn UpdateFn
}

func funcComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNFunc,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Func)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			//sd := data.(*Func)
		},
	})
}

type FuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *FuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		funcComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *FuncComponentSystem) SystemName() string { return SNFunc }

// SystemPriority returns the system priority
func (cs *FuncComponentSystem) SystemPriority() int { return 0 }

func (cs *FuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*Func)
			if c.Fn != nil {
				c.Fn(ctx, m.Entity)
			}
		}
	}
}

type BeforeFuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *BeforeFuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		funcComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *BeforeFuncComponentSystem) SystemName() string { return SNFuncBefore }

// SystemPriority returns the system priority
func (cs *BeforeFuncComponentSystem) SystemPriority() int { return math.MinInt32 }

// SystemExec returns the exec func
func (cs *BeforeFuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*Func)
			if c.BeforeFn != nil {
				c.BeforeFn(ctx, m.Entity)
			}
		}
	}
}

type AfterFuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *AfterFuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		funcComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *AfterFuncComponentSystem) SystemName() string { return SNFuncAfter }

// SystemPriority returns the system priority
func (cs *AfterFuncComponentSystem) SystemPriority() int { return math.MinInt32 }

// SystemExec returns the exec func
func (cs *AfterFuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*Func)
			if c.AfterFn != nil {
				c.AfterFn(ctx, m.Entity)
			}
		}
	}
}

func init() {
	RegisterComponentSystem(&BeforeFuncComponentSystem{})
	RegisterComponentSystem(&FuncComponentSystem{})
	RegisterComponentSystem(&AfterFuncComponentSystem{})
}
