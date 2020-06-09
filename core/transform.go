package core

import (
	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	SNTransform       = "primen.TransformSystem"
	SNTransformSprite = "primen.TransformSpriteSystem"
	CNTransform       = "primen.TransformComponent"
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

func (cs *TransformComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		transformComponentDef(w),
	}
}

func (cs *TransformComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

func transformComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNTransform,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Transform)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
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

	// priv
	lastTick uint64
	// calculated transform matrix (Ebiten)
	m      ebiten.GeoM
	mAngle float64
}

// NewTransform returns a new transform with ScaleX = 1 and ScaleY = 1
func NewTransform() *Transform {
	return &Transform{
		ScaleX: 1,
		ScaleY: 1,
	}
}

type TransformDrawableComponentSystem struct {
	BaseComponentSystem
}

func (cs *TransformDrawableComponentSystem) SystemName() string {
	return SNTransformSprite
}

func (cs *TransformDrawableComponentSystem) SystemPriority() int {
	return -6
}

// SystemInit returns the system init
func (cs *TransformDrawableComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		sys.View().SetOnEntityRemoved(func(e ecs.Entity, w *ecs.World) {
			if getter := w.Component(CNDrawable); getter != nil {
				if vi := getter.Data(e); vi != nil {
					if v, ok := vi.(Drawable); ok {
						v.ClearTransformMatrix()
					}
				}
			}
		})
	}
}

func (cs *TransformDrawableComponentSystem) SystemExec() SystemExecFn {
	return TransformSpriteSystemExec
}

func (cs *TransformDrawableComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		transformComponentDef(w),
		drawableComponentDef(w),
	}
}

func (cs *TransformDrawableComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
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
		_ = resolveTransformM2(t, tick)
	}
}

// TransformSpriteSystemExec is the main function of the TransformSpriteSystem
func TransformSpriteSystemExec(ctx Context) {
	// dt float64, v *ecs.View, s *ecs.System
	v := ctx.System().View()
	matches := v.Matches()
	transformgetter := ctx.World().Component(CNTransform)
	drawablegetter := ctx.World().Component(CNDrawable)
	for _, m := range matches {
		t := m.Components[transformgetter].(*Transform)
		// transform is already resolved because the TransformSystem executed first
		d := m.Components[drawablegetter].(Drawable)
		d.SetTransformMatrix(GeoM2(t.m))
	}
}

func resolveTransformM2(t *Transform, tick uint64) ebiten.GeoM {
	if t == nil {
		return ebiten.GeoM{}
	}
	if t.lastTick == tick {
		return t.m
	}

	parent := resolveTransformM2(t.Parent, tick)
	xb := &ebiten.GeoM{}

	xb.Scale(t.ScaleX, t.ScaleY)
	xb.Rotate(t.Angle)
	xb.Translate(t.X, t.Y)
	xb.Concat(parent)
	t.m = *xb
	t.lastTick = tick
	return t.m
}
func resolveTransformM(t *Transform, tick uint64) (ebiten.GeoM, float64) {
	if t == nil {
		return ebiten.GeoM{}, 0
	}
	if t.lastTick == tick {
		return t.m, t.mAngle
	}

	base, pangle := resolveTransformM(t.Parent, tick)
	xb := &base

	xb.Rotate(pangle)
	xb.Scale(t.ScaleX, t.ScaleY)
	xb.Rotate(t.Angle)
	xb.Translate(t.X, t.Y)
	t.m = *xb
	t.lastTick = tick
	t.mAngle = pangle + t.Angle
	return t.m, t.mAngle
}

func init() {
	RegisterComponentSystem(&TransformComponentSystem{})
	RegisterComponentSystem(&TransformDrawableComponentSystem{})
}
