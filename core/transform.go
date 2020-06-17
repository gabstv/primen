package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// Transform is a hierarchy based matrix
type Transform struct {
	Parent ecs.Entity
	X      float64
	Y      float64
	Angle  float64
	ScaleX float64
	ScaleY float64

	// priv
	lastTick    uint64
	lastParent  ecs.Entity
	lastParentT *Transform
	// calculated transform matrix (Ebiten)
	m ebiten.GeoM
}

//go:generate ecsgen -n Transform -p core -o transform_component.go --component-tpl --vars "UUID=45E8849D-7EA9-4CDC-8AB1-86DB8705C253"

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
		if v.Transform.lastParent == e {
			v.Transform.lastParentT = nil
		}
	}
}

func (s *TransformSystem) setupTransforms() {
	s.tick = 0
}

func (s *TransformSystem) DrawPriority(ctx DrawCtx) {
	// noop
}

func (s *TransformSystem) Draw(ctx DrawCtx) {
	// noop
}

func (s *TransformSystem) UpdatePriority(ctx UpdateCtx) {

}

func (s *TransformSystem) Update(ctx UpdateCtx) {
	tick := s.tick
	s.tick++

	for _, v := range s.V().Matches() {
		if v.Transform.lastParent != v.Transform.Parent {
			if v.Transform.Parent == 0 {
				v.Transform.lastParentT = nil
			} else {
				vd, ok := s.V().Fetch(v.Transform.Parent)
				if !ok {
					// invalid entity passed!
					v.Transform.lastParentT = nil
					v.Transform.Parent = 0
				} else {
					v.Transform.lastParentT = vd.Transform
				}
			}
		}
		v.Transform.lastParent = v.Transform.Parent
	}

	//
	for _, v := range s.V().Matches() {
		_ = resolveTransform(v.Transform, tick)
	}
}

type transformCache struct {
	M ebiten.GeoM
}

func resolveCache(t *Transform)

func resolveTransform(t *Transform, tick uint64) ebiten.GeoM {
	if t == nil {
		return ebiten.GeoM{}
	}
	if t.lastTick == tick {
		return t.m
	}
	parent := resolveTransform(t.lastParentT, tick)
	t.m = ebiten.GeoM{}
	t.m.Scale(t.ScaleX, t.ScaleY)
	t.m.Rotate(t.Angle)
	t.m.Translate(t.X, t.Y)
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

var resizematchDrawableTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetTransformComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *DrawableTransformSystem) DrawPriority(ctx DrawCtx) {
	// noop
}

func (s *DrawableTransformSystem) Draw(ctx DrawCtx) {
	// noop
}

func (s *DrawableTransformSystem) UpdatePriority(ctx UpdateCtx) {

}

func (s *DrawableTransformSystem) Update(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.concatm = v.Transform.m
		v.Drawable.concatset = true
	}
}
