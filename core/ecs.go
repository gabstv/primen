package core

import (
	"github.com/gabstv/ecs/v2"
)

type GameWorld struct {
	*ecs.World
	e        Engine
	disabled bool
}

func (w *GameWorld) Engine() Engine {
	return w.e
}

// SetEnabled is very useful to prevent weird behavior if you're using a goroutine
// to create entities (and components) inside a scene loader.
func (w *GameWorld) SetEnabled(enabled bool) {
	w.disabled = !enabled
}

// Enabled gets if this world is enabled or not
func (w *GameWorld) Enabled() bool {
	return !w.disabled
}

func NewWorld(e Engine) *GameWorld {
	return &GameWorld{
		World: ecs.NewWorld().(*ecs.World),
		e:     e,
	}
}

type World interface {
	ecs.BaseWorld
	Engine() Engine
	SetEnabled(enabled bool)
	Enabled() bool
}

type System interface {
	DrawPriority(ctx DrawCtx)
	Draw(ctx DrawCtx)
	UpdatePriority(ctx UpdateCtx)
	Update(ctx UpdateCtx)
}
