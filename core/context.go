package core

import (
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
)

type Context interface {
	Frame() int64
	DT() float64
	TPS() float64
	Engine() Engine
}

// UpdateCtx is the context passed to every system update function.
type UpdateCtx interface {
	Context
}

// DrawCtx is the context passed to every system update function.
type DrawCtx interface {
	Context
	Renderer() DrawManager
}

type Engine interface {
	AddModule(module Module, priority int)
	NewWorld(priority int) World
	NewWorldWithDefaults(priority int) World
	RemoveWorld(w World)
	Run() error
	Ready() <-chan struct{}
	UpdateFrame() int64
	DrawFrame() int64
	Width() int
	Height() int
	SizeVec() geom.Vec
	AddEventListener(eventName string, fn EventFn) EventID
	RemoveEventListener(id EventID) bool
	DispatchEvent(eventName string, data interface{})
	SetDebugTPS(v bool)
	SetScreenScale(scale float64)
}

type ctxt struct {
	frame  int64
	dt     float64
	tps    float64
	r      DrawManager
	engine Engine
}

func (c *ctxt) Frame() int64 {
	return c.frame
}

func (c *ctxt) DT() float64 {
	return c.dt
}

func (c *ctxt) TPS() float64 {
	return c.tps
}

func (c *ctxt) Renderer() DrawManager {
	return c.r
}

func (c *ctxt) Engine() Engine {
	return c.engine
}

func NewUpdateCtx(e Engine, frame int64, dt, tps float64) UpdateCtx {
	return &ctxt{
		frame:  frame,
		dt:     dt,
		tps:    tps,
		engine: e,
	}
}

func NewDrawCtx(e Engine, frame int64, dt, tps float64, screen *ebiten.Image, rt ...DrawTarget) DrawCtx {
	return &ctxt{
		frame:  frame,
		dt:     dt,
		tps:    tps,
		r:      newDrawManager(screen, rt...),
		engine: e,
	}
}
