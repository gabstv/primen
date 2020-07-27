package components

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
)

//FIXME: review

type UpdateFn func(ctx core.UpdateCtx, e ecs.Entity)
type DrawFn func(ctx core.DrawCtx, e ecs.Entity)

type Function struct {
	DrawPriority   DrawFn
	Draw           DrawFn
	UpdatePriority UpdateFn
	Update         UpdateFn
}

//go:generate ecsgen -n Function -p components -o function_component.go --component-tpl --vars "UUID=C1A2F07B-6EB2-4F83-B20B-0138073786BA"

//go:generate ecsgen -n Function -p components -o function_system.go --system-tpl --vars "Priority=-1000" --vars "UUID=72828866-8D03-4073-82C6-D127A6633521" --components "Function"

var matchFunctionSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetFunctionComponent(w).Flag())
}

var resizematchFunctionSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetFunctionComponent(w).Flag())
}

func (s *FunctionSystem) DrawPriority(ctx core.DrawCtx) {
	for _, v := range s.V().Matches() {
		if v.Function.DrawPriority != nil {
			v.Function.DrawPriority(ctx, v.Entity)
		}
	}
}

func (s *FunctionSystem) Draw(ctx core.DrawCtx) {
	for _, v := range s.V().Matches() {
		if v.Function.Draw != nil {
			v.Function.Draw(ctx, v.Entity)
		}
	}
}

func (s *FunctionSystem) UpdatePriority(ctx core.UpdateCtx) {
	for _, v := range s.V().Matches() {
		if v.Function.UpdatePriority != nil {
			v.Function.UpdatePriority(ctx, v.Entity)
		}
	}
}

func (s *FunctionSystem) Update(ctx core.UpdateCtx) {
	for _, v := range s.V().Matches() {
		if v.Function.Update != nil {
			v.Function.Update(ctx, v.Entity)
		}
	}
}
