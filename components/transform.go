package components

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/geom"
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
}

func NewTransform(x, y float64) Transform {
	return Transform{
		x:      x,
		y:      y,
		scaleX: 1,
		scaleY: 1,
	}
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

func (t *Transform) SetPos(p geom.Vec) *Transform {
	t.x = p.X
	t.y = p.Y
	return t
}

func (t *Transform) X() float64 {
	return t.x
}

func (t *Transform) Y() float64 {
	return t.y
}

func (t *Transform) Pos() geom.Vec {
	return geom.Vec{t.x, t.y}
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

func (t *Transform) GeoM() ebiten.GeoM {
	return t.m
}

//go:generate ecsgen -n Transform -p components -o transform_component.go --component-tpl --vars "UUID=45E8849D-7EA9-4CDC-8AB1-86DB8705C253" --vars "OnAdd=c.setupTransform(e)" --vars "OnResize=c.resized()" --vars "OnWillResize=c.willresize()" --vars "OnRemove=c.removed(e)"

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

//go:generate ecsgen -n Transform -p components -o transform_system.go --system-tpl --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "Setup=s.setupTransforms()" --vars "Priority=100" --vars "UUID=486FA1E8-4C45-48F2-AD8A-02D84C4426C9" --components "Transform" --members "tick=uint64"

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
	if e == 0 {
		// special case where the transform would be nil (if the obj has no parent)
		return gx, gy, true
	}
	ts, ok := s.V().Fetch(e)
	if !ok {
		return 0, 0, false
	}

	pm := ts.Transform.m
	pm.Invert()
	x, y = pm.Apply(gx, gy)

	return x, y, true
}

func (s *TransformSystem) GlobalToLocalTr(gx, gy float64, tr *Transform) (x, y float64) {
	if tr == nil {
		return gx, gy
	}
	pm := tr.m
	pm.Invert()
	x, y = pm.Apply(gx, gy)

	return x, y
}

func (s *TransformSystem) LocalToGlobal(lx, ly float64, e ecs.Entity) (x, y float64, ok bool) {
	ts, ok := s.V().Fetch(e)
	if !ok {
		return 0, 0, false
	}

	pm := ts.Transform.m
	//pm.Invert()
	x, y = pm.Apply(lx, ly)

	return x, y, true
}

func (s *TransformSystem) LocalToGlobalTr(lx, ly float64, tr *Transform) (x, y float64) {
	pm := tr.m
	//pm.Invert()
	x, y = pm.Apply(lx, ly)

	return x, y
}

// DrawPriority noop
func (s *TransformSystem) DrawPriority(ctx core.DrawCtx) {
	if !debug.Draw {
		return
	}
	for _, v := range s.V().Matches() {
		debug.LineM(ctx.Renderer().Screen(), v.Transform.m, -4, 0, 4, 0, debug.PivotColor)
		debug.LineM(ctx.Renderer().Screen(), v.Transform.m, 0, -4, 0, 4, debug.PivotColor)
	}
}

// Draw noop
func (s *TransformSystem) Draw(ctx core.DrawCtx) {}

// UpdatePriority noop
func (s *TransformSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update calculates all transform matrices
func (s *TransformSystem) Update(ctx core.UpdateCtx) {
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
