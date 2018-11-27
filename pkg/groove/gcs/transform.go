package gcs

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/common"
)

const (
	TransformPriority       int = 0
	TransformSpritePriority int = -6
)

var (
	transformWC = &common.WorldComponents{}
)

func init() {
	groove.DefaultComp(func(e *groove.Engine, w *ecs.World) {
		TransformComponent(w)
	})
	groove.DefaultSys(func(e *groove.Engine, w *ecs.World) {
		TransformSystem(w)
		TransformSpriteSystem(w)
	})
	println("transforminit end")
}

// Transform is a hierarchy based matrix
type Transform struct {
	Parent *Transform
	X      float64
	Y      float64
	Angle  float64

	// calculated transform matrix
	M Matrix

	// priv
	lastTick    uint64
	globalAngle float64
}

// TransformComponent will get the registered transform component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func TransformComponent(w *ecs.World) *ecs.Component {
	c := transformWC.Get(w)
	if c == nil {
		var err error
		c, err = w.NewComponent(ecs.NewComponentInput{
			Name: "groove.gcs.Transform",
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*Transform)
				return ok
			},
			DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
				//sd := data.(*Transform)
			},
		})
		if err != nil {
			panic(err)
		}
		transformWC.Set(w, c)
	}
	return c
}

// TransformSystem creates the transform system
func TransformSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(TransformPriority, TransformSystemExec, transformWC.Get(w))
	sys.AddTag(groove.WorldTagUpdate)
	sys.Set("tick", uint64(0))
	return sys
}

// TransformSpriteSystem creates the transform sprite system
func TransformSpriteSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(TransformSpritePriority, TransformSpriteSystemExec, transformWC.Get(w), spriteWC.Get(w))
	sys.AddTag(groove.WorldTagUpdate)
	println("TransformSpriteSystem")
	return sys
}

// TransformSystemExec is the main function of the TransformSystem
func TransformSystemExec(dt float64, v *ecs.View, s *ecs.System) {
	tick := s.Get("tick").(uint64)
	tick++
	s.Set("tick", tick)
	//
	world := v.World()
	matches := v.Matches()
	transformcomp := transformWC.Get(world)
	//engine := world.Get(groove.EngineKey).(*groove.Engine)
	for _, m := range matches {
		t := m.Components[transformcomp].(*Transform)
		resolveTransform(t, tick)
	}
}

// TransformSystemExec is the main function of the TransformSpriteSystem
func TransformSpriteSystemExec(dt float64, v *ecs.View, s *ecs.System) {
	matches := v.Matches()
	world := v.World()
	transformcomp := transformWC.Get(world)
	spritecomp := spriteWC.Get(world)
	for _, m := range matches {
		t := m.Components[transformcomp].(*Transform)
		// transform is already resolved because the TransformSystem executed first
		s := m.Components[spritecomp].(*Sprite)
		vvec := t.M.Project(ZV)
		s.X = vvec.X
		s.Y = vvec.Y
		s.Angle = t.globalAngle
	}
}

func resolveTransform(t *Transform, tick uint64) {
	if t.Parent != nil && t.Parent.lastTick != tick {
		resolveTransform(t.Parent, tick)
	}
	m0 := IM
	pA := float64(0)
	if t.Parent != nil {
		m0 = t.Parent.M
		pA = t.globalAngle
	}

	m1 := m0.Chained(IM)
	if t.X != 0 || t.Y != 0 {
		m1 = m1.Moved(V(t.X, t.Y))
	}
	if t.Angle != 0 {
		m1 = m1.Rotated(ZV, t.Angle)
	}
	t.globalAngle = pA + t.Angle
	t.M = m1
	t.lastTick = tick
}
