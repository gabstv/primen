package tau

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau/core"
)

type Transform struct {
	*WorldItem
	TauTransform *core.Transform
}

func NewTransform(w *ecs.World, parent *core.Transform) *Transform {
	tr := &Transform{
		WorldItem: newWorldItem(w.NewEntity(), w),
	}
	tr.TauTransform = &core.Transform{
		Parent: parent,
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(tr.entity, w.Component(core.CNTransform), tr.TauTransform); err != nil {
		panic(err)
	}
	return tr
}
