package primen

import (
	"sync"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen/core"
)

// Object is the base of any Primen base ECS object
type Object interface {
	Entity() ecs.Entity
	World() *ecs.World
}

// ObjectContainer is an object that contains other objects
type ObjectContainer interface {
	Children() []Object
}

// TransformGetter is the interface to get the core.Transform of an object
type TransformGetter interface {
	CoreTransform() *core.Transform
}

// WorldTransform is minimum interface to most objects with a transform component
type WorldTransform interface {
	TransformGetter
	World() *ecs.World
}

// TransformChild is an interface for objects with a settable parent transform
type TransformChild interface {
	SetParent(parent TransformGetter)
}

// Transformer is the interface of an object with a Transform component
type Transformer interface {
	TransformGetter
	TransformChild
}

// Layerer is the interface of an object with a core.DrawLayer
type Layerer interface {
	SetLayer(l Layer)
	SetZIndex(index int64)
	Layer() Layer
	ZIndex() int64
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

// WorldItem implements Object
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

// Entity an ID
func (wi *WorldItem) Entity() ecs.Entity {
	return wi.entity
}

// World of this object instance
func (wi *WorldItem) World() *ecs.World {
	return wi.world
}

// UpsertFns upserts the core.UpdateFn component to this object's entity
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

// TransformItem implements Transformer
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
		t.transform.Parent = parent.CoreTransform()
	}
	if err := parent.World().AddComponentToEntity(e, parent.World().Component(core.CNTransform), t.transform); err != nil {
		panic(err)
	}
	return t
}

// CoreTransform retrieves the lower level core.Transform
func (t *TransformItem) CoreTransform() *core.Transform {
	return t.transform
}

// SetParent sets the parent transform
func (t *TransformItem) SetParent(parent TransformGetter) {
	t.transform.Parent = parent.CoreTransform()
}

// SetPos sets the transform x and y position (relative to the parent)
func (t *TransformItem) SetPos(x, y float64) {
	t.transform.X = x
	t.transform.Y = y
}

// SetX sets the X position of the transform
func (t *TransformItem) SetX(x float64) {
	t.transform.X = x
}

// SetY sets the Y position of the transform
func (t *TransformItem) SetY(y float64) {
	t.transform.Y = y
}

// SetScale sets the x and y scale of the transform (1 = 100%; 0.0 = 0%)
func (t *TransformItem) SetScale(sx, sy float64) {
	t.transform.ScaleX = sx
	t.transform.ScaleY = sy
}

// SetScale2 sets the x and y scale of the transform (1 = 100%; 0.0 = 0%)
func (t *TransformItem) SetScale2(s float64) {
	t.transform.ScaleX = s
	t.transform.ScaleY = s
}

// SetAngle sets the local angle (in radians) of the transform
func (t *TransformItem) SetAngle(radians float64) {
	t.transform.Angle = radians
}

type engineWT struct {
	w *ecs.World
}

func (wt *engineWT) CoreTransform() *core.Transform {
	return nil
}

func (wt *engineWT) World() *ecs.World {
	return wt.w
}

// Root transform of the world is the world wrapped in a WorldTransform interface
func (e *Engine) Root(w *ecs.World) WorldTransform {
	if w == nil {
		w = e.Default()
	}
	return &engineWT{
		w: w,
	}
}

type DrawLayerItem struct {
	drawLayer *core.DrawLayer
}

func newDrawLayerItem(e ecs.Entity, w *ecs.World) *DrawLayerItem {
	l := &DrawLayerItem{
		drawLayer: &core.DrawLayer{},
	}
	if err := w.AddComponentToEntity(e, w.Component(core.CNTransform), l.drawLayer); err != nil {
		panic(err)
	}
	return l
}

func (dli *DrawLayerItem) SetLayer(l Layer) {
	dli.drawLayer.Layer = l
}
func (dli *DrawLayerItem) SetZIndex(index int64) {
	dli.drawLayer.ZIndex = index
}
func (dli *DrawLayerItem) Layer() Layer {
	return dli.drawLayer.Layer
}
func (dli *DrawLayerItem) ZIndex() int64 {
	return dli.drawLayer.ZIndex
}
