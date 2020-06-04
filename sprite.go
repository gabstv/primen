package primen

import (
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	*WorldItem
	*TransformItem
	CoreSprite *core.Sprite
	DrawLayer  *core.DrawLayer
}

func NewSprite(parent WorldTransform, im *ebiten.Image, layer Layer) *Sprite {
	w := parent.World()
	e := w.NewEntity()
	spr := &Sprite{}
	spr.WorldItem = newWorldItem(e, w)
	spr.TransformItem = newTransformItem(e, parent)
	spr.CoreSprite = &core.Sprite{
		ScaleX: 1,
		ScaleY: 1,
		Bounds: im.Bounds(),
		Image:  im,
	}
	spr.DrawLayer = &core.DrawLayer{
		Layer:  layer,
		ZIndex: core.ZIndexTop,
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawable), spr.CoreSprite); err != nil {
		panic(err)
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawLayer), spr.DrawLayer); err != nil {
		panic(err)
	}
	return spr
}
