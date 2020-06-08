package primen

import (
	"errors"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen/core"
)

// SetFuncs adds a component with the specified update functions
func SetFuncs(w *ecs.World, entity ecs.Entity, beforefn, fn, afterfn core.UpdateFn) error {
	if w == nil {
		return errors.New("world is nil")
	}
	return w.AddComponentToEntity(entity, w.Component(core.CNFunc), &core.Func{
		BeforeFn: beforefn,
		Fn:       fn,
		AfterFn:  afterfn,
	})
}

// SetDrawFuncs adds a component with the specified draw functions
func SetDrawFuncs(w *ecs.World, entity ecs.Entity, beforefn, fn, afterfn core.DrawFn) error {
	if w == nil {
		return errors.New("world is nil")
	}
	return w.AddComponentToEntity(entity, w.Component(core.CNDrawFunc), &core.DrawFunc{
		BeforeFn: beforefn,
		Fn:       fn,
		AfterFn:  afterfn,
	})
}
