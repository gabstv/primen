package primen

import (
	"image/color"

	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
)

var (
	transparentPixel *ebiten.Image
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

type Animation = core.Animation

type AnimatedSprite struct {
	*Sprite
	coreAnim *core.SpriteAnimation
}

func NewAnimatedSprite(parent WorldTransform, layer Layer, anim Animation) *AnimatedSprite {
	as := &AnimatedSprite{
		Sprite: NewSprite(parent, transparentPixel, layer),
	}
	sa := &core.SpriteAnimation{
		Enabled: true,
		Anim:    anim,
	}
	if err := as.World().AddComponentToEntity(as.Entity(), as.World().Component(core.CNSpriteAnimation), sa); err != nil {
		panic(err)
	}
	as.coreAnim = sa
	return as
}

func (as *AnimatedSprite) PlayClipIndex(i int) {
	as.coreAnim.PlayClipIndex(i)
}

func (as *AnimatedSprite) PlayClip(name string) {
	as.coreAnim.PlayClip(name)
}

func init() {
	transparentPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = transparentPixel.Fill(color.Transparent)
}
