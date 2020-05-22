package primen

import (
	"errors"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen/core"
)

// SetFuncs adds a component with the specified functions
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
