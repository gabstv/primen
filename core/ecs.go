package core

import (
	"github.com/gabstv/ecs/v2"
)

type GameWorld struct {
	*ecs.World
	e Engine
}

func (w *GameWorld) Engine() Engine {
	return w.e
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
}

type System interface {
	DrawPriority(ctx DrawCtx)
	Draw(ctx DrawCtx)
	UpdatePriority(ctx UpdateCtx)
	Update(ctx UpdateCtx)
}
