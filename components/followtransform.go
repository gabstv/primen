package components

import (
	"math"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/geom"
)

// we choose not to track the transform data since this component is not very common
// so "targetTr *Transform" is not used. we get from targetE every frame

type FollowTransform struct {
	targetE      ecs.Entity
	withRotation bool
	deadZone     geom.Vec
	bounds       geom.Rect
	offset       geom.Vec
	localOffset  bool
	targetP      geom.Vec
}

func (c *FollowTransform) SetTarget(e ecs.Entity) {
	c.targetE = e
}

func (c *FollowTransform) SetDeadZone(dz geom.Vec) {
	c.deadZone = dz
}

func (c *FollowTransform) SetBounds(dz geom.Rect) {
	c.bounds = dz
}

func (c *FollowTransform) SetOffset(dz geom.Rect) {
	c.bounds = dz
}

func (c *FollowTransform) SetIsLocalOffset(islocal bool) {
	c.localOffset = islocal
}

//go:generate ecsgen -n FollowTransform -p components -o followtransform_component.go --component-tpl --vars "UUID=1245B961-718E-4580-AEB7-893877FD948C"

//go:generate ecsgen -n FollowTransform -p components -o followtransform_transformsystem.go --system-tpl --vars "Priority=60" --vars "UUID=F06A8AE9-8615-4526-A266-E726F4BE0D8A" --components "Transform" --components "FollowTransform"

var matchFollowTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetTransformComponent(w).Flag().Or(GetFollowTransformComponent(w).Flag()))
}

var resizematchFollowTransformSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetTransformComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetFollowTransformComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *FollowTransformSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw noop
func (s *FollowTransformSystem) Draw(ctx core.DrawCtx) {
	// ts := GetTransformSystem(s.world)
	// for _, v := range s.V().Matches() {
	// 	dt := ctx.Renderer().DrawTarget(v.Camera.drawTarget)
	// 	if dt == nil {
	// 		continue
	// 	}
	// 	dt.ResetTransform()
	// 	//TODO: dt.Scale()
	// 	//TODO: dt.Rotate()
	// 	gx, gy := ts.LocalToGlobalTr(0, 0, v.Transform)
	// 	// dt.Translate(geom.Vec{v.Transform.X(), v.Transform.Y()})
	// 	dt.Translate(geom.Vec{-gx, -gy}.Add(v.Camera.viewRect.Scaled(.5)).Add(v.Camera.offset))
	// }
	// if !debug.Draw {
	// 	return
	// }
	// // #ebd951
	// boundsC := color.RGBA{
	// 	R: 0xeb,
	// 	G: 0xd9,
	// 	B: 0x51,
	// 	A: 200,
	// }
	// for _, v := range s.V().Matches() {
	// 	if !v.Camera.deadZone.IsZero() {
	// 		x1, y1 := v.Camera.deadZone.Min.X, v.Camera.deadZone.Min.Y
	// 		x2, y2 := v.Camera.deadZone.Max.X, v.Camera.deadZone.Max.Y
	// 		debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y1, x2, y1, boundsC)
	// 		debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y1, x2, y2, boundsC)
	// 		debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y2, x1, y2, boundsC)
	// 		debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y2, x1, y1, boundsC)
	// 	}
	// }
}

// UpdatePriority noop
func (s *FollowTransformSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update calculates all transform matrices
func (s *FollowTransformSystem) Update(ctx core.UpdateCtx) {
	ts := GetTransformSystem(s.world)
	dt := ctx.DT()

	for _, v := range s.V().Matches() {
		targetTr := GetTransformComponentData(s.world, v.FollowTransform.targetE)
		if targetTr == nil {
			continue
		}
		// get target position on local space
		lx, ly := 0.0, 0.0
		if v.FollowTransform.localOffset {
			lx, ly = v.FollowTransform.offset.X, v.FollowTransform.offset.Y
		}
		gx, gy := ts.LocalToGlobalTr(lx, ly, targetTr)
		lx, ly = ts.GlobalToLocalTr(gx, gy, v.Transform.parent)
		if !v.FollowTransform.localOffset {
			lx += v.FollowTransform.offset.X
			ly += v.FollowTransform.offset.Y
		}
		lv := geom.Vec{lx, ly}
		clampedx := false
		clampedy := false
		if !v.FollowTransform.bounds.IsZero() {
			prev := lv
			lv = lv.RectClamp(v.FollowTransform.bounds)
			if !geom.ScalarEqualsEpsilon(lv.X, prev.X, geom.Epsilon) {
				clampedx = true
			}
			if !geom.ScalarEqualsEpsilon(lv.Y, prev.Y, geom.Epsilon) {
				clampedy = true
			}
		}
		// don't move if inside deadZone
		//FIXME: deadzone code needs work
		if !v.FollowTransform.deadZone.IsZero() {
			dx := lv.X - v.Transform.x
			dy := lv.Y - v.Transform.y
			if math.Abs(dx) < v.FollowTransform.deadZone.X && math.Abs(dy) < v.FollowTransform.deadZone.Y {
				if clampedx && clampedy {
					v.FollowTransform.targetP = lv
				} else if clampedx {
					v.FollowTransform.targetP = geom.Vec{lv.X, v.FollowTransform.targetP.Y}
				} else if clampedy {
					v.FollowTransform.targetP = geom.Vec{v.FollowTransform.targetP.X, lv.Y}
				}
				continue
			}
			// needs to move, but take deadZone offset (hw, hh) into consideration
			if !clampedx {
				if lv.X > v.Transform.x {
					lv.X -= v.FollowTransform.deadZone.X / 2.0 //-math.Abs(dx)
				} else {
					lv.X += v.FollowTransform.deadZone.X / 2.0
				}
			}
			if !clampedy {
				if lv.Y > v.Transform.y {
					lv.Y -= v.FollowTransform.deadZone.Y / 2.0 //-math.Abs(dY)
				} else {
					lv.Y += v.FollowTransform.deadZone.Y / 2.0
				}
			}
			lv = lv.RectClamp(v.FollowTransform.bounds)
		}
		v.FollowTransform.targetP = lv
	}
	for _, v := range s.V().Matches() {
		lp := geom.Vec{v.Transform.x, v.Transform.y}
		ft := defaultRTEasing(lp, v.FollowTransform.targetP, dt)
		v.Transform.x = ft.X
		v.Transform.y = ft.Y
	}
	// // tick := s.tick
	// // s.tick++
	// //
	// ts := GetTransformSystem(s.world)
	// dt := ctx.DT()
	// for _, v := range s.V().Matches() {

	// 	targetTr := GetTransformComponentData(s.world, v.Camera.targetE)
	// 	if targetTr == nil {
	// 		continue
	// 	}
	// 	gx, gy := ts.LocalToGlobalTr(0, 0, targetTr)
	// 	lx, ly := ts.GlobalToLocalTr(gx, gy, v.Transform.parent)
	// 	//TODO: use deadzone
	// 	if !v.Camera.deadZone.IsZero() {
	// 		if v.Camera.insideDeadZone(ts, v.Transform, geom.Vec{lx, ly}) {
	// 			continue
	// 		}
	// 	}
	// 	v.Camera.lastTargetPos = geom.Vec{lx, ly}
	// }
	// for _, v := range s.V().Matches() {
	// 	// movement happens here
	// 	mvpos := v.Camera.ease(geom.Vec{v.Transform.X(), v.Transform.Y()}, dt)
	// 	v.Transform.SetX(mvpos.X).SetY(mvpos.Y)
	// }
}

func defaultRTEasing(curp, targetp geom.Vec, dt float64) geom.Vec {
	spdv := targetp.Sub(curp)
	//mag := spdv.Magnitude()
	addv := spdv.Scaled(2 * dt)
	return curp.Add(addv)
}
