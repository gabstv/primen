package primen

import "github.com/gabstv/primen/core"

type Animation = core.Animation

type AnimatedSprite struct {
	*Sprite
	coreAnim func() *core.SpriteAnimation
}

func NewAnimatedSprite(parent WorldTransform, layer Layer, anim Animation) *AnimatedSprite {
	as := &AnimatedSprite{
		Sprite: NewSprite(parent, transparentPixel, layer),
	}
	core.SetSpriteAnimationComponentData(as.world, as.entity, core.SpriteAnimation{
		Enabled: true,
		Anim:    anim,
	})
	as.coreAnim = func() *core.SpriteAnimation { return core.GetSpriteAnimationComponentData(as.world, as.entity) }
	return as
}

func (as *AnimatedSprite) PlayClipIndex(i int) {
	as.coreAnim().PlayClipIndex(i)
}

func (as *AnimatedSprite) PlayClip(name string) {
	as.coreAnim().PlayClip(name)
}
