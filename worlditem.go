package primen

import (
	"sync"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen/core"
)

type Object interface {
	Entity() ecs.Entity
	World() *ecs.World
}

type ObjectContainer interface {
	Children() []Object
}

// Destroy an object
func Destroy(obj Object) bool {
	if v, ok := obj.(ObjectContainer); ok {
		for _, child := range v.Children() {
			_ = Destroy(child)
		}
	}
	return obj.World().RemoveEntity(obj.Entity())
}

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

type WorldTransform interface {
	TransformGetter
	World() *ecs.World
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
	children  []Object
	childrenm sync.Mutex
}

func newTransformItem(e ecs.Entity, parent WorldTransform) *TransformItem {
	if parent == nil {
		panic("parent can't be nil. Use e.Root(nil) if this object need to be at root")
	}
	t := &TransformItem{
		transform: core.NewTransform(),
	}
	if parent != nil {
		t.transform.Parent = parent.GetCoreTransform()
	}
	if err := parent.World().AddComponentToEntity(e, parent.World().Component(core.CNTransform), t.transform); err != nil {
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

type engineWT struct {
	w *ecs.World
}

func (wt *engineWT) GetCoreTransform() *core.Transform {
	return nil
}

func (wt *engineWT) World() *ecs.World {
	return wt.w
}

func (e *Engine) Root(w *ecs.World) WorldTransform {
	if w == nil {
		w = e.Default()
	}
	return &engineWT{
		w: w,
	}
}
