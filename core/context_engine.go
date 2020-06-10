package core

import (
	"context"
	"time"

	"github.com/gabstv/ecs"
)

// Context is the context passed to every system update function.
type Context interface {
	ecs.Context
	Engine() Engine
	FPS() float64
	Frame() int64
	IsDrawingSkipped() bool
	Renderer() DrawManager
}

type Engine interface {
	AddWorld(w *ecs.World, priority int)
	RemoveWorld(w *ecs.World) bool
	Default() *ecs.World
	Run() error
	Ready() <-chan struct{}
	Get(key string) interface{}
	Set(key string, value interface{})
	UpdateFrame() int64
	DrawFrame() int64
	Width() int
	Height() int
	AddEventListener(eventName string, fn EventFn) EventID
	RemoveEventListener(id EventID) bool
	DispatchEvent(eventName string, data interface{})
}

type ctxt struct {
	c          context.Context
	dt         float64
	system     *ecs.System
	world      *ecs.World
	engine     Engine
	fps        float64
	frame      int64
	drwskipped bool
	drawM      DrawManager
}

func (c ctxt) Deadline() (deadline time.Time, ok bool) {
	return c.c.Deadline()
}

func (c ctxt) Done() <-chan struct{} {
	return c.c.Done()
}

func (c ctxt) Err() error {
	return c.c.Err()
}

func (c ctxt) Value(key interface{}) interface{} {
	return c.c.Value(key)
}

func (c ctxt) DT() float64 {
	return c.dt
}

func (c ctxt) System() *ecs.System {
	return c.system
}

func (c ctxt) World() *ecs.World {
	return c.world
}

func (c ctxt) Engine() Engine {
	return c.engine
}

func (c ctxt) FPS() float64 {
	return c.fps
}

func (c ctxt) Frame() int64 {
	return c.frame
}

func (c ctxt) IsDrawingSkipped() bool {
	return c.drwskipped
}

func (c ctxt) Renderer() DrawManager {
	return c.drawM
}
