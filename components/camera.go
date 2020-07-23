package components

import (
	"image/color"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
)

type CameraRealTimeEasingFn func(curp, targetp geom.Vec, dt float64) geom.Vec

// we choose not to track the transform data since cameras are not very common
// so "targetTr *Transform" is not used. we get from targetE every frame

type Camera struct {
	targetE       ecs.Entity
	withRotation  bool
	viewRect      geom.Vec
	deadZone      geom.Rect
	offset        geom.Vec
	rteasing      CameraRealTimeEasingFn
	inverted      bool
	drawTarget    core.DrawTargetID
	lastTargetPos geom.Vec
}

func NewCamera(drawTarget core.DrawTargetID) Camera {
	return Camera{
		drawTarget: drawTarget,
	}
}

func (c *Camera) SetViewRect(r geom.Vec) {
	c.viewRect = r
}

func (c *Camera) SetDeadZone(dz geom.Rect) {
	c.deadZone = dz
}

func (c *Camera) SetTarget(e ecs.Entity) {
	c.targetE = e
}

func (c *Camera) SetInverted(inverted bool) {
	c.inverted = inverted
}

func (c *Camera) ease(curp geom.Vec, dt float64) geom.Vec {
	if c.rteasing == nil {
		return defaultCameraEasing(curp, c.lastTargetPos, dt)
	}
	return c.rteasing(curp, c.lastTargetPos, dt)
}

//go:generate ecsgen -n Camera -p components -o camera_component.go --component-tpl --vars "UUID=02167024-81F4-4FEB-AC18-36564FCAC20B"

//go:generate ecsgen -n Camera -p components -o camera_transformsystem.go --system-tpl --vars "Priority=50" --vars "UUID=E19B710D-139B-47BD-AF0C-340414BC7226" --components "Transform" --components "Camera"

var matchCameraSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	return f.Contains(GetTransformComponent(w).Flag().Or(GetCameraComponent(w).Flag()))
}

var resizematchCameraSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetTransformComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetCameraComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *CameraSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw noop
func (s *CameraSystem) Draw(ctx core.DrawCtx) {
	ts := GetTransformSystem(s.world)
	for _, v := range s.V().Matches() {
		dt := ctx.Renderer().DrawTarget(v.Camera.drawTarget)
		if dt == nil {
			continue
		}
		dt.ResetTransform()
		//TODO: dt.Scale()
		//TODO: dt.Rotate()
		gx, gy := ts.LocalToGlobalTr(0, 0, v.Transform)
		// dt.Translate(geom.Vec{v.Transform.X(), v.Transform.Y()})
		dt.Translate(geom.Vec{-gx, -gy}.Add(v.Camera.viewRect.Scaled(.5)))
	}
	if !debug.Draw {
		return
	}
	// #ebd951
	boundsC := color.RGBA{
		R: 0xeb,
		G: 0xd9,
		B: 0x51,
		A: 200,
	}
	for _, v := range s.V().Matches() {
		if !v.Camera.deadZone.IsZero() {
			x1, y1 := v.Camera.deadZone.Min.X, v.Camera.deadZone.Min.Y
			x2, y2 := v.Camera.deadZone.Max.X, v.Camera.deadZone.Max.Y
			debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y1, x2, y1, boundsC)
			debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y1, x2, y2, boundsC)
			debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y2, x1, y2, boundsC)
			debug.LineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y2, x1, y1, boundsC)
		}
	}
}

// UpdatePriority noop
func (s *CameraSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update calculates all transform matrices
func (s *CameraSystem) Update(ctx core.UpdateCtx) {
	// tick := s.tick
	// s.tick++
	//
	ts := GetTransformSystem(s.world)
	dt := ctx.DT()
	for _, v := range s.V().Matches() {

		targetTr := GetTransformComponentData(s.world, v.Camera.targetE)
		if targetTr == nil {
			//TODO: move by easing and final position etc
			mvpos := v.Camera.ease(geom.Vec{v.Transform.X(), v.Transform.Y()}, dt)
			v.Transform.SetX(mvpos.X).SetY(mvpos.Y)
			continue
		}
		gx, gy := ts.LocalToGlobalTr(0, 0, targetTr)
		lx, ly := ts.GlobalToLocalTr(gx, gy, v.Transform.parent)
		//TODO: use deadzone
		v.Camera.lastTargetPos = geom.Vec{lx, ly}
		//
		mvpos := v.Camera.ease(geom.Vec{v.Transform.X(), v.Transform.Y()}, dt)
		v.Transform.SetX(mvpos.X).SetY(mvpos.Y)
		//
	}
}

func defaultCameraEasing(curp, targetp geom.Vec, dt float64) geom.Vec {
	spdv := targetp.Sub(curp)
	//mag := spdv.Magnitude()
	addv := spdv.Scaled(.5)
	return curp.Add(addv)
}
