package gcs

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/common"
	"github.com/hajimehoshi/ebiten"
)

const (
	TransformPriority int = 0
)

var (
	transformWC *common.WorldComponents
)

func init() {
	transformWC = &common.WorldComponents{}
	//
	groove.DefaultComp(func(e *groove.Engine, w *ecs.World) {
		TransformComponent(w)
	})
	//groove.DefaultSys(func(e *groove.Engine, w *ecs.World) {
	//	TransformSystem(w)
	//})
}

// Transform is a hierarchy based matrix
type Transform struct {
	Parent *Transform
	X      float64
	Y      float64
	Angle  float64

	// calculated transform matrix
	M ebiten.GeoM

	// priv
	lastTick uint64
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

func resolveTransform(t *Transform, tick uint64) {
	if t.Parent != nil && t.Parent.lastTick != tick {
		resolveTransform(t.Parent, tick)
	}
	var m0 ebiten.GeoM
	if t.Parent != nil {
		m0 = t.Parent.M
	}
	//TODO: test
	m1 := t.M
	m1.Reset()
	m1.Concat(m0)
	m1.Translate(t.X, t.Y)
	m1.Rotate(t.Angle)
	t.M = m1
	t.lastTick = tick
}
