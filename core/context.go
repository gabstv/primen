package core

import "github.com/hajimehoshi/ebiten"

//FIXME: review

type Context interface {
	Frame() int64
	DT() float64
	TPS() float64
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
	NewWorld(priority int) World
	NewWorldWithDefaults(priority int) World
	RemoveWorld(w World)
	Run() error
	Ready() <-chan struct{}
	UpdateFrame() int64
	DrawFrame() int64
	Width() int
	Height() int
	AddEventListener(eventName string, fn EventFn) EventID
	RemoveEventListener(id EventID) bool
	DispatchEvent(eventName string, data interface{})
	SetDebugTPS(v bool)
}

type ctxt struct {
	frame int64
	dt    float64
	tps   float64
	r     DrawManager
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

func NewUpdateCtx(frame int64, dt, tps float64) UpdateCtx {
	return &ctxt{
		frame: frame,
		dt:    dt,
		tps:   tps,
	}
}

func NewDrawCtx(frame int64, dt, tps float64, screen *ebiten.Image) DrawCtx {
	return &ctxt{
		frame: frame,
		dt:    dt,
		tps:   tps,
		r:     newDrawManager(screen),
	}
}

// type ctxt struct {
// 	c          context.Context
// 	dt         float64
// 	system     *ecs.System
// 	world      *ecs.World
// 	engine     Engine
// 	fps        float64
// 	frame      int64
// 	drwskipped bool
// 	drawM      DrawManager
// }

// func (c ctxt) Deadline() (deadline time.Time, ok bool) {
// 	return c.c.Deadline()
// }

// func (c ctxt) Done() <-chan struct{} {
// 	return c.c.Done()
// }

// func (c ctxt) Err() error {
// 	return c.c.Err()
// }

// func (c ctxt) Value(key interface{}) interface{} {
// 	return c.c.Value(key)
// }

// func (c ctxt) DT() float64 {
// 	return c.dt
// }

// func (c ctxt) System() *ecs.System {
// 	return c.system
// }

// func (c ctxt) World() *ecs.World {
// 	return c.world
// }

// func (c ctxt) Engine() Engine {
// 	return c.engine
// }

// func (c ctxt) FPS() float64 {
// 	return c.fps
// }

// func (c ctxt) Frame() int64 {
// 	return c.frame
// }

// func (c ctxt) IsDrawingSkipped() bool {
// 	return c.drwskipped
// }

// func (c ctxt) Renderer() DrawManager {
// 	return c.drawM
// }
