package core

import (
	"github.com/gabstv/ecs/v2"
)

//FIXME: review

type GameWorld struct {
	*ecs.World
	e Engine
}

func NewWorld(e Engine) *GameWorld {
	return &GameWorld{
		World: ecs.NewWorld().(*ecs.World),
		e:     e,
	}
}

type World interface {
	ecs.BaseWorld
	DrawPriority(ctx DrawCtx)
	Draw(ctx DrawCtx)
	UpdatePriority(ctx UpdateCtx)
	Update(ctx UpdateCtx)
}

type System interface {
	DrawPriority(ctx DrawCtx)
	Draw(ctx DrawCtx)
	UpdatePriority(ctx UpdateCtx)
	Update(ctx UpdateCtx)
}
