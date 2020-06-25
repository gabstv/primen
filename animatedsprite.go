package primen

import (
	"github.com/gabstv/primen/core"
)

type Animation = core.Animation

type AnimatedSpriteNode struct {
	*SpriteNode
	ca core.WatchSpriteAnimation
}

func NewRootAnimatedSpriteNode(w World, layer Layer, fps float64, anim core.Animation) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewRootSpriteNode(w, layer),
	}
	core.SetSpriteAnimationComponentData(w, sprn.e, core.NewSpriteAnimation(fps, anim))
	sprn.ca = core.WatchSpriteAnimationComponentData(w, sprn.e)
	return sprn
}

func NewChildAnimatedSpriteNode(parent ObjectContainer, layer Layer, fps float64, anim core.Animation) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewChildSpriteNode(parent, layer),
	}
	core.SetSpriteAnimationComponentData(parent.World(), sprn.e, core.NewSpriteAnimation(fps, anim))
	sprn.ca = core.WatchSpriteAnimationComponentData(parent.World(), sprn.e)
	return sprn
}

func (n *AnimatedSpriteNode) SpriteAnim() *core.SpriteAnimation {
	return n.ca.Data()
}
