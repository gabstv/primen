package primen

import "github.com/gabstv/primen/core"

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
