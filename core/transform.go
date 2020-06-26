package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// Transform is a hierarchy based matrix
type Transform struct {
	x      float64
	y      float64
	angle  float64
	scaleX float64
	scaleY float64

	// priv
	lastTick uint64
	pentity  ecs.Entity
	parent   *Transform
	// calculated transform matrix (Ebiten)
	m ebiten.GeoM
	// copy the current world
	// this is set once the component is created
	w ecs.BaseWorld
	// called before removal
	removenotice map[int64]func()
	removeindex  int64
}

func NewTransform(x, y float64) Transform {
	return Transform{
		x:            x,
		y:            y,
		scaleX:       1,
		scaleY:       1,
		removenotice: make(map[int64]func()),
	}
}

func (t *Transform) AddDestroyListener(l func()) int64 {
	t.removeindex++
	id := t.removeindex
	t.removenotice[id] = l
	return id
}
func (t *Transform) RemoveDestroyListener(id int64) {
	delete(t.removenotice, id)
}

func (t *Transform) SetParent(e ecs.Entity) bool {
	pt := GetTransformComponentData(t.w, e)
	if pt == nil {
		return false
	}
	//TODO: check cyclic transform parenting
	t.pentity = e
	t.parent = pt
	return true
}

func (t *Transform) Parent() ecs.Entity {
	return t.pentity
}

func (t *Transform) ParentTransform() *Transform {
	return t.parent
}

func (t *Transform) Tree() []*Transform {
	tr := make([]*Transform, 0, 16)
	t.tree(&tr)
	return tr
}

func (t *Transform) tree(l *[]*Transform) {
	*l = append(*l, t)
	if t.parent != nil {
		t.parent.tree(l)
	}
}

func (t *Transform) SetX(x float64) *Transform {
	t.x = x
	return t
}

func (t *Transform) SetY(y float64) *Transform {
	t.y = y
	return t
}

func (t *Transform) X() float64 {
	return t.x
}

func (t *Transform) Y() float64 {
	return t.y
}

// SetAngle sets the angle (in radians)
func (t *Transform) SetAngle(r float64) *Transform {
	t.angle = r
	return t
}

// Angle gets the angle (in radians)
func (t *Transform) Angle() float64 {
	return t.angle
}

func (t *Transform) SetScale(sx, sy float64) *Transform {
	t.scaleX, t.scaleY = sx, sy
	return t
}

func (t *Transform) Scale() (sx, sy float64) {
	return t.scaleX, t.scaleY
}

func (t *Transform) SetScaleX(sx float64) *Transform {
	t.scaleX = sx
	return t
}

func (t *Transform) ScaleX() float64 {
	return t.scaleX
}

func (t *Transform) SetScaleY(sy float64) *Transform {
	t.scaleY = sy
	return t
}

func (t *Transform) ScaleY() float64 {
	return t.scaleY
}

//go:generate ecsgen -n Transform -p core -o transform_component.go --component-tpl --vars "UUID=45E8849D-7EA9-4CDC-8AB1-86DB8705C253" --vars "OnAdd=c.setupTransform(e)" --vars "OnResize=c.resized()" --vars "OnWillResize=c.willresize()" --vars "OnRemove=c.removed(e)" --vars "BeforeRemove=c.beforeremove(e)"

func (c *TransformComponent) setupTransform(e ecs.Entity) {
	d := c.Data(e)
	d.w = c.world
}

func (c *TransformComponent) willresize() {
	for i, v := range c.data {
		d := v.Data
		d.parent = nil
		v.Data = d
		c.data[i] = v
	}
}

func (c *TransformComponent) resized() {
	for i, v := range c.data {
		if v.Data.pentity == 0 {
			v.Data.parent = nil
		} else {
			x := &v.Data
			if d := c.Data(x.pentity); d != nil {
				x.parent = d
			} else {
				x.pentity = 0
				x.parent = nil
			}
		}
		c.data[i] = v
	}
}

func (c *TransformComponent) removed(e ecs.Entity) {
	for i := range c.data {
		if c.data[i].Data.pentity == e {
			x := &c.data[i].Data
			x.parent = nil
			x.pentity = 0
		}
	}
}

func (c *TransformComponent) beforeremove(e ecs.Entity) {
	i := c.indexof(e)
	if c.data[i].Data.removenotice != nil {
		for _, v := range c.data[i].Data.removenotice {
			v()
		}
	}
	c.data[i].Data.removenotice = nil
}

//go:generate ecsgen -n Transform -p core -o transform_system.go --system-tpl --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "Setup=s.setupTransforms()" --vars "Priority=100" --vars "UUID=486FA1E8-4C45-48F2-AD8A-02D84C4426C9" --components "Transform" --members "tick=uint64"

var matchTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetTransformComponent(w).Flag())
}

var resizematchTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetTransformComponent(w).Flag())
}

func (s *TransformSystem) onEntityAdded(e ecs.Entity) {

}

func (s *TransformSystem) onEntityRemoved(e ecs.Entity) {
	for _, v := range s.V().Matches() {
		if v.Transform.pentity == e {
			v.Transform.parent = nil
		}
	}
}

func (s *TransformSystem) setupTransforms() {
	s.tick = 0
}

func (s *TransformSystem) GlobalToLocal(gx, gy float64, e ecs.Entity) (x, y float64, ok bool) {
	ts, ok := s.V().Fetch(e)
	if !ok {
		return 0, 0, false
	}
	//ts.Transform.Parent()
	//m := ebiten.GeoM{}
	//m.Translate(gx, gy)
	//m2 := ts.Transform.m
	//m.Apply()
	//x, y = ts.Transform.m.Apply(gx, gy)
	//return x, y, true
	// m := ebiten.GeoM{}
	// m.Translate(gx, gy)
	// m.Invert()
	// m2 := ts.Transform.m
	// m2.Concat(m)
	// m2.Invert()
	// x, y = m2.Apply(0, 0)

	// M_loc = M_parent_inv * M
	pm := ts.Transform.m
	pm.Invert()
	m := ebiten.GeoM{}
	m.Translate(gx, gy)
	pm.Concat(m)
	x, y = pm.Apply(0, 0)
	return x, y, true
}

// DrawPriority noop
func (s *TransformSystem) DrawPriority(ctx DrawCtx) {
	if !DebugDraw {
		return
	}
	for _, v := range s.V().Matches() {
		debugLineM(ctx.Renderer().Screen(), v.Transform.m, -4, 0, 4, 0, debugPivotColor)
		debugLineM(ctx.Renderer().Screen(), v.Transform.m, 0, -4, 0, 4, debugPivotColor)
	}
}

// Draw noop
func (s *TransformSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *TransformSystem) UpdatePriority(ctx UpdateCtx) {}

// Update calculates all transform matrices
func (s *TransformSystem) Update(ctx UpdateCtx) {
	tick := s.tick
	s.tick++
	//
	for _, v := range s.V().Matches() {
		_ = resolveTransform(v.Transform, tick)
	}
}

type transformCache struct {
	M ebiten.GeoM
}

func resolveTransform(t *Transform, tick uint64) ebiten.GeoM {
	if t == nil {
		return ebiten.GeoM{}
	}
	if t.lastTick == tick {
		return t.m
	}
	parent := resolveTransform(t.parent, tick)
	t.m = ebiten.GeoM{}
	t.m.Scale(t.scaleX, t.scaleY)
	t.m.Rotate(t.angle)
	t.m.Translate(t.x, t.y)
	t.m.Concat(parent)
	t.lastTick = tick
	return t.m
}

//go:generate ecsgen -n DrawableTransform -p core -o transform_drawablesystem.go --system-tpl --vars "Priority=90" --vars "UUID=7E9DEBA9-DEF6-4174-8160-AA7B72E2A734" --components "Transform" --components "Drawable"

var matchDrawableTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetTransformComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	return true
}

// The DrawableTransformSystem's View needs to be recalculated if the
// Drawable or Transform component arrays change.
var resizematchDrawableTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetTransformComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *DrawableTransformSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop
func (s *DrawableTransformSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *DrawableTransformSystem) UpdatePriority(ctx UpdateCtx) {}

// Update sets the drawable transform
func (s *DrawableTransformSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.SetConcatM(v.Transform.m)
	}
}
