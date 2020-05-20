package tau

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau/core"
	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	*WorldItem
	TauSprite *core.Sprite
	Transform *core.Transform
	DrawLayer *core.DrawLayer
}

func NewSprite(w *ecs.World, im *ebiten.Image, layer core.LayerIndex, parent *core.Transform) *Sprite {
	spr := &Sprite{}
	spr.WorldItem = newWorldItem(w.NewEntity(), w)
	spr.TauSprite = &core.Sprite{
		ScaleX: 1,
		ScaleY: 1,
		Bounds: im.Bounds(),
		Image:  im,
	}
	spr.DrawLayer = &core.DrawLayer{
		Layer:  layer,
		ZIndex: core.ZIndexTop,
	}
	spr.Transform = &core.Transform{
		Parent: parent,
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawable), spr.TauSprite); err != nil {
		panic(err)
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawLayer), spr.DrawLayer); err != nil {
		panic(err)
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNTransform), spr.Transform); err != nil {
		panic(err)
	}
	return spr
}
