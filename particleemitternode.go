package primen

import (
	"github.com/gabstv/primen/components/graphics"
)

type ParticleEmitterNode struct {
	*Node
	wdl graphics.WatchDrawLayer
	wpe graphics.WatchParticleEmitter
}

func NewRootParticleEmitterNode(w World, layer Layer) *ParticleEmitterNode {
	pen := &ParticleEmitterNode{
		Node: NewRootNode(w),
	}
	graphics.SetDrawLayerComponentData(w, pen.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetParticleEmitterComponentData(w, pen.e, graphics.NewParticleEmitter(w))
	pen.wdl = graphics.WatchDrawLayerComponentData(w, pen.e)
	pen.wpe = graphics.WatchParticleEmitterComponentData(w, pen.e)
	return pen
}

func NewChildParticleEmitterNode(parent ObjectContainer, layer Layer) *ParticleEmitterNode {
	pen := &ParticleEmitterNode{
		Node: NewChildNode(parent),
	}
	graphics.SetDrawLayerComponentData(parent.World(), pen.e, graphics.DrawLayer{
		Layer:  layer,
		ZIndex: graphics.ZIndexTop,
	})
	graphics.SetParticleEmitterComponentData(parent.World(), pen.e, graphics.NewParticleEmitter(parent.World()))
	pen.wdl = graphics.WatchDrawLayerComponentData(parent.World(), pen.e)
	pen.wpe = graphics.WatchParticleEmitterComponentData(parent.World(), pen.e)
	return pen
}

func (n *ParticleEmitterNode) ParticleEmitter() *graphics.ParticleEmitter {
	return n.wpe.Data()
}

func (n *ParticleEmitterNode) SetLayer(l Layer) {
	n.wdl.Data().Layer = l
}

func (n *ParticleEmitterNode) SetZIndex(index int64) {
	n.wdl.Data().ZIndex = index
}

func (n *ParticleEmitterNode) Layer() Layer {
	return n.wdl.Data().Layer
}

func (n *ParticleEmitterNode) ZIndex() int64 {
	return n.wdl.Data().ZIndex
}
