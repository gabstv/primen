package primen

import (
	// 	"sync"

	"sort"

	"github.com/gabstv/ecs/v2"
)

// Object is the base of any Primen base ECS object
type Object interface {
	Entity() ecs.Entity
	World() World
	Destroy()
	SetParent(parent ObjectContainer)
}

// ObjectContainer is an object that contains other objects
type ObjectContainer interface {
	Object
	Children() []Object
	AddChild(child Object)
	RemoveChild(child Object)
}

type mObject struct {
	e      ecs.Entity
	w      World
	parent ObjectContainer
}

type mObjectContainer struct {
	*mObject
	children []Object
}

func (o *mObject) Entity() ecs.Entity {
	return o.e
}

func (o *mObject) World() World {
	return o.w
}

func (o *mObject) SetParent(parent ObjectContainer) {
	o.parent = parent
}

func (o *mObject) Destroy() {
	if o.parent != nil {
		o.parent.RemoveChild(o)
		o.parent = nil
	}
	if o.w != nil && o.e != 0 {
		o.w.RemoveEntity(o.e)
	}
	o.w = nil
	o.e = 0
}

func (o *mObjectContainer) Children() []Object {
	return o.children
}

func (o *mObjectContainer) Destroy() {
	if o.children != nil {
		clone := make([]Object, len(o.children))
		copy(clone, o.children)
		for _, v := range clone {
			if v != nil {
				v.Destroy()
			}
		}
		o.children = nil
	}
	if o.mObject != nil {
		o.mObject.Destroy()
	}
	o.mObject = nil
}

func (o *mObjectContainer) AddChild(child Object) {
	if child == nil {
		return
	}
	if o.children == nil {
		o.children = make([]Object, 0, 2)
	}
	o.children = append(o.children, child)
	sort.Slice(o.children, func(i, j int) bool {
		return o.children[i].Entity() < o.children[j].Entity()
	})
}

func (o *mObjectContainer) RemoveChild(child Object) {
	if child == nil {
		return
	}
	e := child.Entity()
	i := sort.Search(len(o.children), func(i int) bool { return o.children[i].Entity() >= e })
	if i < len(o.children) && o.children[i].Entity() == e {
		o.children = o.children[:i+copy(o.children[i:], o.children[i+1:])]
	}
}
