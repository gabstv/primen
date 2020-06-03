package core

import (
	"math"

	"github.com/gabstv/ecs"
)

const (
	CNDrawFunc       = "tau/core.DrawFuncComponent"
	SNDrawFunc       = "tau/core.DrawFuncSystem"
	SNDrawFuncAfter  = "tau/core.DrawFuncAfterSystem"
	SNDrawFuncBefore = "tau/core.DrawFuncBeforeSystem"
)

type DrawFn func(ctx Context, e ecs.Entity)

type DrawFunc struct {
	Fn       DrawFn
	AfterFn  DrawFn
	BeforeFn DrawFn
}

func drawFuncComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNDrawFunc,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*DrawFunc)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			//sd := data.(*DrawFunc)
		},
	})
}

type DrawFuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *DrawFuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawFuncComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *DrawFuncComponentSystem) SystemName() string { return SNDrawFunc }

// SystemPriority returns the system priority
func (cs *DrawFuncComponentSystem) SystemPriority() int { return 0 }

func (cs *DrawFuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNDrawFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*DrawFunc)
			if c.Fn != nil {
				c.Fn(ctx, m.Entity)
			}
		}
	}
}

func (cs *DrawFuncComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

type BeforeDrawFuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *BeforeDrawFuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawFuncComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *BeforeDrawFuncComponentSystem) SystemName() string { return SNDrawFuncBefore }

// SystemPriority returns the system priority
func (cs *BeforeDrawFuncComponentSystem) SystemPriority() int { return math.MinInt32 }

// SystemExec returns the exec func
func (cs *BeforeDrawFuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNDrawFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*DrawFunc)
			if c.BeforeFn != nil {
				c.BeforeFn(ctx, m.Entity)
			}
		}
	}
}

func (cs *BeforeDrawFuncComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

type AfterDrawFuncComponentSystem struct {
	BaseComponentSystem
}

// Components returns the component signature(s)
func (cs *AfterDrawFuncComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawFuncComponentDef(w),
	}
}

// SystemName returns the system name
func (cs *AfterDrawFuncComponentSystem) SystemName() string { return SNDrawFuncAfter }

// SystemPriority returns the system priority
func (cs *AfterDrawFuncComponentSystem) SystemPriority() int { return math.MinInt32 }

// SystemExec returns the exec func
func (cs *AfterDrawFuncComponentSystem) SystemExec() SystemExecFn {
	return func(ctx Context) {
		getter := ctx.World().Component(CNDrawFunc)
		for _, m := range ctx.System().View().Matches() {
			c := m.Components[getter].(*DrawFunc)
			if c.AfterFn != nil {
				c.AfterFn(ctx, m.Entity)
			}
		}
	}
}

func (cs *AfterDrawFuncComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

func init() {
	RegisterComponentSystem(&BeforeDrawFuncComponentSystem{})
	RegisterComponentSystem(&DrawFuncComponentSystem{})
	RegisterComponentSystem(&AfterDrawFuncComponentSystem{})
}
