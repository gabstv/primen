package primen

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
)

type TransformGetter interface {
	Entity() ecs.Entity
	Transform() *core.Transform
}

type Node struct {
	*mObjectContainer
	wtr core.WatchTransform
}

func NewRootNode(w World) *Node {
	tr := &Node{
		mObjectContainer: &mObjectContainer{
			mObject: &mObject{
				e: w.NewEntity(),
				w: w,
			},
		},
	}
	core.SetTransformComponentData(w, tr.Entity(), core.NewTransform(0, 0))
	tr.wtr = core.WatchTransformComponentData(w, tr.Entity())
	return tr
}

func NewChildNode(parent ObjectContainer) *Node {
	if parent == nil {
		panic("parent can't be nil")
	}
	tr := &Node{
		mObjectContainer: &mObjectContainer{
			mObject: &mObject{
				e: parent.World().NewEntity(),
				w: parent.World(),
			},
		},
	}
	core.SetTransformComponentData(tr.World(), tr.Entity(), core.NewTransform(0, 0))
	tr.wtr = core.WatchTransformComponentData(tr.World(), tr.Entity())
	tr.SetParent(parent)
	return tr
}

func (t *Node) Transform() *core.Transform {
	return t.wtr.Data()
}

func (t *Node) SetParent(parent ObjectContainer) {
	if t.parent != nil {
		if _, ok := t.parent.(TransformGetter); ok {
			t.wtr.Data().SetParent(0)
		}
		t.parent.RemoveChild(t)
	}
	if parent == nil {
		t.parent = nil
		return
	}
	if p, ok := parent.(TransformGetter); ok {
		t.wtr.Data().SetParent(p.Entity())
	}
	t.mObject.SetParent(parent)
}

func (t *Node) Destroy() {
	t.wtr = nil
	t.mObjectContainer.Destroy()
}

//TODO: implement Detroy

// type Transform struct {
// 	*WorldItem
// 	*TransformItem
// }

// func NewTransform(parent WorldTransform) *Transform {
// 	e := parent.World().NewEntity()
// 	tr := &Transform{
// 		WorldItem:     newWorldItem(e, parent.World()),
// 		TransformItem: newTransformItem(e, parent),
// 	}
// 	return tr
// }
