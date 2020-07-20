package core

import (
	"image/color"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
)

type CameraRealTimeEasingFn func(curp, targetp geom.Vec, dt float64) geom.Vec

type Camera struct {
	targetE ecs.Entity
	// we not to track the transform data since cameras are not very common
	//targetTr *Transform
	withRotation bool
	viewRect     geom.Vec
	deadZone     geom.Rect
	offset       geom.Vec
	rteasing     CameraRealTimeEasingFn
	inverted     bool
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

//go:generate ecsgen -n Camera -p core -o camera_component.go --component-tpl --vars "UUID=02167024-81F4-4FEB-AC18-36564FCAC20B"

//go:generate ecsgen -n Camera -p core -o camera_transformsystem.go --system-tpl --vars "Priority=50" --vars "UUID=E19B710D-139B-47BD-AF0C-340414BC7226" --components "Transform" --components "Camera"

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
func (s *CameraSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop
func (s *CameraSystem) Draw(ctx DrawCtx) {
	if !DebugDraw {
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
			debugLineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y1, x2, y1, boundsC)
			debugLineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y1, x2, y2, boundsC)
			debugLineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x2, y2, x1, y2, boundsC)
			debugLineM(ctx.Renderer().Screen(), ebiten.GeoM{}, x1, y2, x1, y1, boundsC)
		}
	}
}

// UpdatePriority noop
func (s *CameraSystem) UpdatePriority(ctx UpdateCtx) {}

// Update calculates all transform matrices
func (s *CameraSystem) Update(ctx UpdateCtx) {
	// tick := s.tick
	// s.tick++
	//
	ts := GetTransformSystem(s.world)
	dt := ctx.DT()
	for _, v := range s.V().Matches() {
		s.updateCamera(v, ts, dt)
	}
}

func (s *CameraSystem) updateCamera(v VICameraSystem, ts *TransformSystem, dt float64) {
	targetTr := GetTransformComponentData(s.world, v.Camera.targetE)
	if targetTr == nil {
		//TODO: move by easing and final position etc
		return
	}
	gpos := getCameraGlobalPos(v.Camera, targetTr, ts)
	if !v.Camera.deadZone.IsZero() {
		// don't move if inside the dead zone
		x, y := ts.LocalToGlobalTr(0, 0, v.Transform)
		x -= v.Camera.viewRect.X / 2
		y -= v.Camera.viewRect.Y / 2
		curgpos := geom.Vec{x, y}
		ds := v.Camera.deadZone.SubVec(v.Camera.viewRect)
		if ds.ContainsVec(curgpos) &&
			ds.ContainsVec(gpos) {
			// do nothing
			return
		}
	}
	lpos := getCameraLocalPos(gpos, v.Camera, v.Transform, ts)
	if v.Camera.inverted {
		lpos = lpos.Scaled(-1)
	}
	cpos := geom.Vec{v.Transform.x, v.Transform.y}
	if v.Camera.rteasing == nil {
		cpos = defaultCameraEasing(cpos, lpos, dt)
		v.Transform.x, v.Transform.y = cpos.X, cpos.Y
	} else {
		cpos = v.Camera.rteasing(cpos, lpos, dt)
		v.Transform.x, v.Transform.y = cpos.X, cpos.Y
	}
	//TODO: camera rotation
}

func defaultCameraEasing(curp, targetp geom.Vec, dt float64) geom.Vec {
	spdv := targetp.Sub(curp)
	//mag := spdv.Magnitude()
	addv := spdv.Scaled(.5)
	return curp.Add(addv)
}

func getCameraGlobalPos(c *Camera, target *Transform, tsys *TransformSystem) geom.Vec {
	// if p := target.ParentTransform(); p != nil {
	// 	gx, gy := tsys.LocalToGlobalTr(target.x, target.y, p)
	// 	return geom.Vec{gx, gy}
	// }
	gx, gy := tsys.LocalToGlobalTr(0, 0, target)
	gx -= c.viewRect.X / 2
	gy -= c.viewRect.Y / 2
	return geom.Vec{gx, gy}
}

func getCameraLocalPos(gpos geom.Vec, c *Camera, tr *Transform, tsys *TransformSystem) geom.Vec {
	lx, ly := tsys.GlobalToLocalTr(gpos.X, gpos.Y, tr)
	return geom.Vec{lx + c.offset.X, ly + c.offset.Y}
}
