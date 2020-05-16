package graphics

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	Entity    ecs.Entity
	TauSprite *tau.Sprite
	Transform *tau.Transform
	DrawLayer *tau.DrawLayer
}

func NewSprite(w *ecs.World, im *ebiten.Image, layer tau.LayerIndex, parent *tau.Transform) *Sprite {
	spr := &Sprite{}
	spr.Entity = w.NewEntity()
	spr.TauSprite = &tau.Sprite{
		ScaleX: 1,
		ScaleY: 1,
		Bounds: im.Bounds(),
		Image:  im,
	}
	spr.DrawLayer = &tau.DrawLayer{
		Layer:  layer,
		ZIndex: tau.ZIndexTop,
	}
	spr.Transform = &tau.Transform{
		Parent: parent,
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(spr.Entity, w.Component(tau.CNSprite), spr.TauSprite); err != nil {
		panic(err)
	}
	if err := w.AddComponentToEntity(spr.Entity, w.Component(tau.CNDrawLayer), spr.DrawLayer); err != nil {
		panic(err)
	}
	if err := w.AddComponentToEntity(spr.Entity, w.Component(tau.CNTransform), spr.Transform); err != nil {
		panic(err)
	}
	return spr
}
