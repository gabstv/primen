package hello

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const SPEED float64 = 120

type Hello struct {
	Text string
	X    int
	Y    int
}

//go:generate ecsgen -n Hello -p hello -o hello_component.go --component-tpl --vars "UUID=3E674178-4393-423A-9C07-99AD18C0CDD0"

type Move struct {
	XSpeed float64
	YSpeed float64
	XSum   float64
	YSum   float64
}

//go:generate ecsgen -n Move -p hello -o move_component.go --component-tpl --vars "UUID=94091420-5E23-457E-8F8D-A422A03E36AF"

//go:generate ecsgen -n Hello -p hello -o hello_system.go --system-tpl --vars "Priority=0" --vars "UUID=C435CD6B-100E-4440-8224-72497D18C6D8" --components "Hello"

var matchHelloSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(GetHelloComponent(w).Flag())
}

var resizematchHelloSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(GetHelloComponent(w).Flag())
}

// DrawPriority noop
func (s *HelloSystem) DrawPriority(ctx core.DrawCtx) {

}

// Draw text via ebitenutil
func (s *HelloSystem) Draw(ctx core.DrawCtx) {
	for _, v := range s.V().Matches() {
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), v.Hello.Text, v.Hello.X, v.Hello.Y)
	}
}

// UpdatePriority noop
func (s *HelloSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update computes labes if dirty
func (s *HelloSystem) Update(ctx core.UpdateCtx) {}

//go:generate ecsgen -n MoveHello -p hello -o movehello_system.go --system-tpl --vars "Priority=10" --vars "UUID=2E58680B-E13E-4763-A997-C744E4107820" --components "Hello" --components "Move"

var matchMoveHelloSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(GetHelloComponent(w).Flag().Or(GetMoveComponent(w).Flag()))
}

var resizematchMoveHelloSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	if eflag.Contains(GetHelloComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(GetMoveComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *MoveHelloSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw noop
func (s *MoveHelloSystem) Draw(ctx core.DrawCtx) {}

// UpdatePriority noop
func (s *MoveHelloSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update positions
func (s *MoveHelloSystem) Update(ctx core.UpdateCtx) {
	dt := ctx.DT()
	for _, v := range s.V().Matches() {
		v.Move.XSum += dt * v.Move.XSpeed
		v.Move.YSum += dt * v.Move.YSpeed
		for v.Move.XSum >= 1 {
			v.Hello.X++
			v.Move.XSum--
		}
		for v.Move.XSum <= -1 {
			v.Hello.X--
			v.Move.XSum++
		}
		for v.Move.YSum >= 1 {
			v.Hello.Y++
			v.Move.YSum--
		}
		for v.Move.YSum <= -1 {
			v.Hello.Y--
			v.Move.YSum++
		}
		if v.Hello.X >= 280 && v.Move.XSpeed > 0 {
			v.Move.XSpeed = -SPEED
		}
		if v.Hello.X <= 0 && v.Move.XSpeed < 0 {
			v.Move.XSpeed = SPEED
		}
		if v.Hello.Y >= 220 && v.Move.YSpeed > 0 {
			v.Move.YSpeed = -SPEED
		}
		if v.Hello.Y <= 0 && v.Move.YSpeed < 0 {
			v.Move.YSpeed = SPEED
		}
	}
}
