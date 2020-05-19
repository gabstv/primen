package tau

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau/core"
)

type Transform struct {
	Entity       ecs.Entity
	TauTransform *core.Transform
}

func NewTransform(w *ecs.World, parent *core.Transform) *Transform {
	tr := &Transform{
		Entity: w.NewEntity(),
	}
	tr.TauTransform = &core.Transform{
		Parent: parent,
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(tr.Entity, w.Component(core.CNTransform), tr.TauTransform); err != nil {
		panic(err)
	}
	return tr
}
