package core

import (
	"github.com/gabstv/ecs"
)

const (
	SNRotationTransform = "tau.RotationTransformSystem"
	CNRotation          = "tau.RotationComponent"
)

type RotationTransformComponentSystem struct {
	BaseComponentSystem
}

func (cs *RotationTransformComponentSystem) SystemName() string {
	return SNRotationTransform
}

func (cs *RotationTransformComponentSystem) SystemPriority() int {
	return 1
}

func (cs *RotationTransformComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		sys.Set("tick", uint64(0))
	}
}

func (cs *RotationTransformComponentSystem) SystemExec() SystemExecFn {
	return RotationTransformSystemExec
}

func (cs *RotationTransformComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		transformComponentDef(w),
		rotationComponentDef(w),
	}
}

func rotationComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNRotation,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Rotation)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			//sd := data.(*Transform)
		},
	})
}

type Rotation struct {
	Speed float64
}

// RotationTransformSystemExec is the main function of the RotationTransformSystem
func RotationTransformSystemExec(ctx Context) {
	// dt float64, v *ecs.View, s *ecs.System
	s := ctx.System()
	v := s.View()
	dt := ctx.DT()
	//
	matches := v.Matches()
	tgetter := ctx.World().Component(CNTransform)
	rgetter := ctx.World().Component(CNRotation)
	for _, m := range matches {
		t := m.Components[tgetter].(*Transform)
		r := m.Components[rgetter].(*Rotation)
		t.Angle += r.Speed * dt
	}
}

func init() {
	RegisterComponentSystem(&RotationTransformComponentSystem{})
}
