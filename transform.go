package tau

import (
	"github.com/gabstv/ecs"
)

const (
	SNTransform       = "tau.TransformSystem"
	SNTransformSprite = "tau.TransformSpriteSystem"
	CNTransform       = "tau.TransformComponent"
)

type TransformComponentSystem struct {
	BaseComponentSystem
}

func (cs *TransformComponentSystem) SystemName() string {
	return SNTransform
}

func (cs *TransformComponentSystem) SystemPriority() int {
	return 0
}

func (cs *TransformComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		sys.Set("tick", uint64(0))
	}
}

func (cs *TransformComponentSystem) SystemExec() SystemExecFn {
	return TransformSystemExec
}

func (cs *TransformComponentSystem) Components(w ecs.Worlder) []*ecs.Component {
	return []*ecs.Component{
		transformComponentDef(w),
	}
}

func transformComponentDef(w ecs.Worlder) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNTransform,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Transform)
			return ok
		},
		DestructorFn: func(_ ecs.WorldDicter, entity ecs.Entity, data interface{}) {
			//sd := data.(*Transform)
		},
	})
}

// Transform is a hierarchy based matrix
type Transform struct {
	Parent *Transform
	X      float64
	Y      float64
	Angle  float64
	ScaleX float64
	ScaleY float64

	// calculated transform matrix
	M Matrix

	// priv
	lastTick     uint64
	globalAngle  float64
	globalScaleX float64
	globalScaleY float64
}

// NewTransform returns a new transform with ScaleX = 1 and ScaleY = 1
func NewTransform() *Transform {
	return &Transform{
		ScaleX: 1,
		ScaleY: 1,
	}
}

type TransformSpriteComponentSystem struct {
	BaseComponentSystem
}

func (cs *TransformSpriteComponentSystem) SystemName() string {
	return SNTransformSprite
}

func (cs *TransformSpriteComponentSystem) SystemPriority() int {
	return -6
}

func (cs *TransformSpriteComponentSystem) SystemExec() SystemExecFn {
	return TransformSpriteSystemExec
}

func (cs *TransformSpriteComponentSystem) Components(w ecs.Worlder) []*ecs.Component {
	return []*ecs.Component{
		transformComponentDef(w),
		spriteComponentDef(w),
	}
}

// TransformSystemExec is the main function of the TransformSystem
func TransformSystemExec(ctx Context) {
	// dt float64, v *ecs.View, s *ecs.System
	s := ctx.System()
	v := s.View()
	tick := s.Get("tick").(uint64)
	tick++
	s.Set("tick", tick)
	//
	matches := v.Matches()
	transformcomp := ctx.World().Component(CNTransform)
	for _, m := range matches {
		t := m.Components[transformcomp].(*Transform)
		resolveTransform(t, tick)
	}
}

// TransformSpriteSystemExec is the main function of the TransformSpriteSystem
func TransformSpriteSystemExec(ctx Context) {
	// dt float64, v *ecs.View, s *ecs.System
	v := ctx.System().View()
	matches := v.Matches()
	transformcomp := ctx.World().Component(CNTransform)
	spritecomp := ctx.World().Component(CNSprite)
	for _, m := range matches {
		t := m.Components[transformcomp].(*Transform)
		// transform is already resolved because the TransformSystem executed first
		s := m.Components[spritecomp].(*Sprite)
		vvec := t.M.Project(ZV)
		s.X = vvec.X
		s.Y = vvec.Y
		s.Angle = t.globalAngle
		s.ScaleX = t.globalScaleX
		s.ScaleY = t.globalScaleY
		//TODO: convert pixel matrix to ebiten matrix (and use it for scale/skew)
	}
}

func resolveTransform(t *Transform, tick uint64) {
	if t.Parent != nil && t.Parent.lastTick != tick {
		resolveTransform(t.Parent, tick)
	}
	parentAngle := float64(0)
	parentScaleX := float64(1)
	parentScaleY := float64(1)
	localAngle := t.Angle
	parentMatrix := IM
	if t.Parent != nil {
		parentAngle = t.Parent.globalAngle
		parentScaleX = t.Parent.globalScaleX
		parentScaleY = t.Parent.globalScaleY
		parentMatrix = t.Parent.M
	}
	t.M = IM.ScaledXY(ZV, V(t.ScaleX, t.ScaleY)).Rotated(ZV, localAngle).Moved(V(t.X, t.Y)).Chained(parentMatrix)
	t.globalAngle = parentAngle + localAngle
	t.globalScaleX = parentScaleX * t.ScaleX
	t.globalScaleY = parentScaleY * t.ScaleY
	t.lastTick = tick
}

func init() {
	RegisterComponentSystem(&TransformComponentSystem{})
	RegisterComponentSystem(&TransformSpriteComponentSystem{})
}
