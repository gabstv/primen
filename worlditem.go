package primen

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/primen/core"
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

type TransformGetter interface {
	GetCoreTransform() *core.Transform
}

type TransformSetter interface {
	SetParent(parent TransformGetter)
}

type Transformer interface {
	TransformGetter
	TransformSetter
}

type TransformItem struct {
	transform *core.Transform
}

func newTransformItem(e ecs.Entity, w *ecs.World, parent TransformGetter) *TransformItem {
	t := &TransformItem{
		transform: core.NewTransform(),
	}
	if parent != nil {
		t.transform.Parent = parent.GetCoreTransform()
	}
	if err := w.AddComponentToEntity(e, w.Component(core.CNTransform), t.transform); err != nil {
		panic(err)
	}
	return t
}

func (t *TransformItem) GetCoreTransform() *core.Transform {
	return t.transform
}

func (t *TransformItem) SetParent(parent TransformGetter) {
	t.transform.Parent = parent.GetCoreTransform()
}

func (t *TransformItem) SetPos(x, y float64) {
	t.transform.X = x
	t.transform.Y = y
}

func (t *TransformItem) SetX(x float64) {
	t.transform.X = x
}

func (t *TransformItem) SetY(y float64) {
	t.transform.Y = y
}

func (t *TransformItem) SetScale(sx, sy float64) {
	t.transform.ScaleX = sx
	t.transform.ScaleY = sy
}

func (t *TransformItem) SetScale2(s float64) {
	t.transform.ScaleX = s
	t.transform.ScaleY = s
}

func (t *TransformItem) SetAngle(radians float64) {
	t.transform.Angle = radians
}
