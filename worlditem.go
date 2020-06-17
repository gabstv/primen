package primen

import (
	"sync"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
)

// Object is the base of any Primen base ECS object
type Object interface {
	Entity() ecs.Entity
	World() core.World
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
	World() *core.GameWorld
	Entity() ecs.Entity
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
	world  *core.GameWorld
}

func newWorldItem(e ecs.Entity, w *core.GameWorld) *WorldItem {
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
func (wi *WorldItem) World() *core.GameWorld {
	return wi.world
}

// UpsertFns upserts the core.UpdateFn component to this object's entity
func (wi *WorldItem) UpsertFns(drawpriority, draw core.DrawFn, updatepriority, update core.UpdateFn) bool {
	if vi := core.GetFunctionComponentData(wi.world, wi.entity); vi != nil {
		vi.DrawPriority = drawpriority
		vi.Draw = draw
		vi.UpdatePriority = updatepriority
		vi.Update = update
		return false
	}
	core.SetFunctionComponentData(wi.world, wi.entity, core.Function{
		DrawPriority:   drawpriority,
		Draw:           draw,
		UpdatePriority: updatepriority,
		Update:         update,
	})
	return true
}

// TransformItem implements Transformer
type TransformItem struct {
	transform func() *core.Transform
	children  []Object
	childrenm sync.Mutex
}

func newTransformItem(e ecs.Entity, parent WorldTransform) *TransformItem {
	if parent == nil {
		panic("parent can't be nil. Use e.Root(nil) if this object need to be at root")
	}
	t := &TransformItem{
		//transform: core.NewTransform(),
	}
	core.SetTransformComponentData(parent.World(), e, core.Transform{})
	w := parent.World()
	t.transform = func() *core.Transform { return core.GetTransformComponentData(w, e) }
	if parent != nil {
		t.transform().Parent = parent.Entity()
	}
	return t
}

// CoreTransform retrieves the lower level core.Transform
func (t *TransformItem) CoreTransform() *core.Transform {
	return t.transform()
}

// SetParent sets the parent transform
func (t *TransformItem) SetParent(parent ecs.Entity) {
	t.transform().Parent = parent
}

// SetPos sets the transform x and y position (relative to the parent)
func (t *TransformItem) SetPos(x, y float64) {
	t.transform().X = x
	t.transform().Y = y
}

// SetX sets the X position of the transform
func (t *TransformItem) SetX(x float64) {
	t.transform().X = x
}

// X gets the x position of the transform
func (t *TransformItem) X() float64 {
	return t.transform().X
}

// Y gets the y position of the transform
func (t *TransformItem) Y() float64 {
	return t.transform().Y
}

// SetY sets the Y position of the transform
func (t *TransformItem) SetY(y float64) {
	t.transform().Y = y
}

// SetScale sets the x and y scale of the transform (1 = 100%; 0.0 = 0%)
func (t *TransformItem) SetScale(sx, sy float64) {
	t.transform().ScaleX = sx
	t.transform().ScaleY = sy
}

// SetScale2 sets the x and y scale of the transform (1 = 100%; 0.0 = 0%)
func (t *TransformItem) SetScale2(s float64) {
	t.transform().ScaleX = s
	t.transform().ScaleY = s
}

func (t *TransformItem) Scale() (x, y float64) {
	return t.transform().ScaleX, t.transform().ScaleY
}

// SetAngle sets the local angle (in radians) of the transform
func (t *TransformItem) SetAngle(radians float64) {
	t.transform().Angle = radians
}

// Angle gets the local angle (in radians) of the transform
func (t *TransformItem) Angle() (radians float64) {
	return t.transform().Angle
}

type engineWT struct {
	w *core.GameWorld
}

func (wt *engineWT) CoreTransform() *core.Transform {
	return nil
}

func (wt *engineWT) World() *core.GameWorld {
	return wt.w
}

func (wt *engineWT) Entity() ecs.Entity {
	return 0
}

// Root transform of the world is the world wrapped in a WorldTransform interface
func (e *Engine) Root(w *core.GameWorld) WorldTransform {
	if w == nil {
		w = e.Default()
	}
	return &engineWT{
		w: w,
	}
}

type DrawLayerItem struct {
	drawLayer func() *core.DrawLayer
}

func newDrawLayerItem(e ecs.Entity, w *core.GameWorld) *DrawLayerItem {
	l := &DrawLayerItem{
		//drawLayer: &core.DrawLayer{},
	}
	core.SetDrawLayerComponentData(w, e, core.DrawLayer{})
	l.drawLayer = func() *core.DrawLayer { return core.GetDrawLayerComponentData(w, e) }
	return l
}

func (dli *DrawLayerItem) SetLayer(l Layer) {
	dli.drawLayer().Layer = l
}
func (dli *DrawLayerItem) SetZIndex(index int64) {
	dli.drawLayer().ZIndex = index
}
func (dli *DrawLayerItem) Layer() Layer {
	return dli.drawLayer().Layer
}
func (dli *DrawLayerItem) ZIndex() int64 {
	return dli.drawLayer().ZIndex
}
