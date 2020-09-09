package components

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/geom"
)

type CameraRealTimeEasingFn func(curp, targetp geom.Vec, dt float64) geom.Vec

// we choose not to track the transform data since cameras are not very common
// so "targetTr *Transform" is not used. we get from targetE every frame

type Camera struct {
	viewRect   geom.Vec
	drawTarget core.DrawTargetID
}

func NewCamera(drawTarget core.DrawTargetID) Camera {
	return Camera{
		drawTarget: drawTarget,
	}
}

func (c *Camera) SetViewRect(r geom.Vec) {
	c.viewRect = r
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
}

// UpdatePriority noop
func (s *CameraSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update noop
func (s *CameraSystem) Update(ctx core.UpdateCtx) {}
