package primen

import (
	"github.com/gabstv/primen/core"
)

type Animation = core.Animation

type AnimatedSpriteNode struct {
	*SpriteNode
	ca core.WatchSpriteAnimation
}

func NewRootAnimatedSpriteNode(w World, layer Layer) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewRootSpriteNode(w, layer),
	}
	core.SetSpriteAnimationComponentData(w, sprn.e, core.SpriteAnimation{})
	sprn.ca = core.WatchSpriteAnimationComponentData(w, sprn.e)
	return sprn
}

func NewChildAnimatedSpriteNode(parent ObjectContainer, layer Layer) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewChildSpriteNode(parent, layer),
	}
	core.SetSpriteAnimationComponentData(parent.World(), sprn.e, core.SpriteAnimation{})
	sprn.ca = core.WatchSpriteAnimationComponentData(parent.World(), sprn.e)
	return sprn
}

//TODO: continue

// func NewAnimatedSprite(parent WorldTransform, layer Layer, anim Animation) *AnimatedSprite {
// 	as := &AnimatedSprite{
// 		Sprite: NewSprite(parent, transparentPixel, layer),
// 	}
// 	core.SetSpriteAnimationComponentData(as.world, as.entity, core.SpriteAnimation{
// 		Enabled: true,
// 		Anim:    anim,
// 	})
// 	as.coreAnim = func() *core.SpriteAnimation { return core.GetSpriteAnimationComponentData(as.world, as.entity) }
// 	return as
// }

// func (as *AnimatedSprite) PlayClipIndex(i int) {
// 	as.coreAnim().PlayClipIndex(i)
// }

// func (as *AnimatedSprite) PlayClip(name string) {
// 	as.coreAnim().PlayClip(name)
// }
