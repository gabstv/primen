package primen

import (
	"github.com/gabstv/primen/core"
)

type ParticleEmitterNode struct {
	*Node
	wdl core.WatchDrawLayer
	wpe core.WatchParticleEmitter
}

func NewRootParticleEmitterNode(w World, layer Layer) *ParticleEmitterNode {
	pen := &ParticleEmitterNode{
		Node: NewRootNode(w),
	}
	core.SetDrawLayerComponentData(w, pen.e, core.DrawLayer{
		Layer:  layer,
		ZIndex: core.ZIndexTop,
	})
	core.SetParticleEmitterComponentData(w, pen.e, core.NewParticleEmitter(w))
	pen.wdl = core.WatchDrawLayerComponentData(w, pen.e)
	pen.wpe = core.WatchParticleEmitterComponentData(w, pen.e)
	return pen
}

func NewChildParticleEmitterNode(parent ObjectContainer, layer Layer) *ParticleEmitterNode {
	pen := &ParticleEmitterNode{
		Node: NewChildNode(parent),
	}
	core.SetDrawLayerComponentData(parent.World(), pen.e, core.DrawLayer{
		Layer:  layer,
		ZIndex: core.ZIndexTop,
	})
	core.SetParticleEmitterComponentData(parent.World(), pen.e, core.NewParticleEmitter(parent.World()))
	pen.wdl = core.WatchDrawLayerComponentData(parent.World(), pen.e)
	pen.wpe = core.WatchParticleEmitterComponentData(parent.World(), pen.e)
	return pen
}

func (n *ParticleEmitterNode) ParticleEmitter() *core.ParticleEmitter {
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
