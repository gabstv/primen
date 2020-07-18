package core

import (
	"github.com/gabstv/primen/geom"
	"github.com/gabstv/ecs/v2"
)

type Camera struct {
	targetE ecs.Entity
	// we not to track the transform data since cameras are not very common
	//targetTr *Transform
	withRotation bool
	deadZone geom.Rect
	offset geom.Vec
}

//go:generate ecsgen -n Camera -p core -o camera_component.go --component-tpl --vars "UUID=02167024-81F4-4FEB-AC18-36564FCAC20B"

//go:generate ecsgen -n Camera -p core -o camera_transformsystem.go --system-tpl --vars "Priority=50" --vars "UUID=6389F54D-76C9-49FC-B3E3-1C73B334EBB6" --components "Transform" --components "Camera"

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
func (s *CameraSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *CameraSystem) UpdatePriority(ctx UpdateCtx) {}

// Update calculates all transform matrices
func (s *CameraSystem) Update(ctx UpdateCtx) {
	tick := s.tick
	s.tick++
	//
	ts := GetTransformSystem(s.world)
	for _, v := range s.V().Matches() {
		targetTr := GetTransformComponentData(s.world, v.Camera.targetE)
		targetTr.ParentTransform()
		ts.LocalToGlobalTr().
	}
}
