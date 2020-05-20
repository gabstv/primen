package core

import (
	"context"
	"time"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

// Context is the context passed to every system update function.
type Context interface {
	ecs.Context
	Engine() Engine
	FPS() float64
	Frame() int64
	IsDrawingSkipped() bool
	DefaultDrawImageOptions() *ebiten.DrawImageOptions
	Screen() *ebiten.Image
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
	imopt      *ebiten.DrawImageOptions
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

func (c ctxt) DefaultDrawImageOptions() *ebiten.DrawImageOptions {
	return c.imopt
}

func (c ctxt) Screen() *ebiten.Image {
	return c.world.Get("screen").(*ebiten.Image)
}
