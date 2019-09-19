package common

import (
	"github.com/gabstv/ecs"
)

// Archetype is a recipe to create entities with a preset of components.
type Archetype struct {
	World      *ecs.World
	Components []*ecs.Component
}

// NewArchetype returns a new archetype. This func unsures that no
// duplicated components are added.
func NewArchetype(world *ecs.World, comps ...*ecs.Component) *Archetype {
	cmap := make(map[*ecs.Component]bool)
	components := make([]*ecs.Component, 0, len(comps))
	for _, c := range comps {
		if cmap[c] {
			// duplicated component
			continue
		}
		components = append(components, c)
		cmap[c] = true
	}
	arch := &Archetype{
		World:      world,
		Components: components,
	}
	return arch
}

// NewEntity adds a new entity with the component data to the world
//
// The best way to create an archetype entity is to ensure that the
// component data follows the same order that the components were
// created in NewArchetype
func (a *Archetype) NewEntity(compdata ...interface{}) ecs.Entity {
	entity := a.World.NewEntity()
	compds := clonecompslc(a.Components)
	for _, cdata := range compdata {
		for i, c := range compds {
			if c.Validate(cdata) {
				if err := a.World.AddComponentToEntity(entity, c, cdata); err != nil {
					// this should never happen
					panic(err)
				}
				compds = append(compds[:i], compds[i+1:]...)
				break
			}
		}
	}
	return entity
}

func clonecompslc(v []*ecs.Component) []*ecs.Component {
	x := make([]*ecs.Component, len(v), len(v))
	for i, vv := range v {
		x[i] = vv
	}
	return x
}
