package groove

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove/common"
)

const (
	// TransformPriority - default 0
	TransformPriority int = 0
	// TransformSpritePriority - default -6 (execs after positioning the transforms)
	TransformSpritePriority int = -6
)

var (
	transformWC = &common.WorldComponents{}
)

func init() {
	DefaultComp(func(e *Engine, w *ecs.World) {
		TransformComponent(w)
	})
	DefaultSys(func(e *Engine, w *ecs.World) {
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

// TransformComponent will get the registered transform component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func TransformComponent(w *ecs.World) *ecs.Component {
	c := transformWC.Get(w)
	if c == nil {
		var err error
		c, err = w.NewComponent(ecs.NewComponentInput{
			Name: "groove.Transform",
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
	sys.AddTag(WorldTagUpdate)
	sys.Set("tick", uint64(0))
	return sys
}

// TransformSpriteSystem creates the transform sprite system
func TransformSpriteSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(TransformSpritePriority, TransformSpriteSystemExec, transformWC.Get(w), spriteWC.Get(w))
	sys.AddTag(WorldTagUpdate)
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
	for _, m := range matches {
		t := m.Components[transformcomp].(*Transform)
		resolveTransform(t, tick)
	}
}

// TransformSpriteSystemExec is the main function of the TransformSpriteSystem
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
