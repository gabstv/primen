package primen

import (
	"github.com/gabstv/primen/components/graphics"
)

type Animation = graphics.Animation

type AnimatedSpriteNode struct {
	*SpriteNode
	ca graphics.WatchSpriteAnimation
}

func NewRootAnimatedSpriteNode(w World, layer Layer, fps float64, anim graphics.Animation) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewRootSpriteNode(w, layer),
	}
	graphics.SetSpriteAnimationComponentData(w, sprn.e, graphics.NewSpriteAnimation(fps, anim))
	sprn.ca = graphics.WatchSpriteAnimationComponentData(w, sprn.e)
	return sprn
}

func NewChildAnimatedSpriteNode(parent ObjectContainer, layer Layer, fps float64, anim graphics.Animation) *AnimatedSpriteNode {
	sprn := &AnimatedSpriteNode{
		SpriteNode: NewChildSpriteNode(parent, layer),
	}
	graphics.SetSpriteAnimationComponentData(parent.World(), sprn.e, graphics.NewSpriteAnimation(fps, anim))
	sprn.ca = graphics.WatchSpriteAnimationComponentData(parent.World(), sprn.e)
	return sprn
}

func (n *AnimatedSpriteNode) SpriteAnim() *graphics.SpriteAnimation {
	return n.ca.Data()
}
