package tau

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau/core"
)

type WorldItem struct {
	entity ecs.Entity
	world  *ecs.World
}

func newWorldItem(e ecs.Entity, w *ecs.World) *WorldItem {
	return &WorldItem{
		entity: e,
		world:  w,
	}
}

func (wi *WorldItem) Entity() ecs.Entity {
	return wi.entity
}

func (wi *WorldItem) World() *ecs.World {
	return wi.world
}

func (wi *WorldItem) UpsertFns(beforefn, fn, afterfn core.UpdateFn) bool {
	if vi := wi.world.Component(core.CNFunc).Data(wi.entity); vi != nil {
		if v, ok := vi.(*core.Func); ok {
			v.BeforeFn = beforefn
			v.Fn = fn
			v.AfterFn = afterfn
		}
		return false
	}
	if err := wi.world.AddComponentToEntity(wi.entity, wi.world.Component(core.CNFunc), &core.Func{
		BeforeFn: beforefn,
		Fn:       fn,
		AfterFn:  afterfn,
	}); err != nil {
		println(err)
		return false
	}
	return true
}
