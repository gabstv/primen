package graphics

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
)

type Transform struct {
	Entity       ecs.Entity
	TauTransform *tau.Transform
}

func NewTransform(w *ecs.World, parent *tau.Transform) *Transform {
	tr := &Transform{
		Entity: w.NewEntity(),
	}
	tr.TauTransform = &tau.Transform{
		Parent: parent,
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(tr.Entity, w.Component(tau.CNTransform), tr.TauTransform); err != nil {
		panic(err)
	}
	return tr
}
