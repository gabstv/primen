package graphics

import (
	"math"
	"sort"
	"sync"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
)

// LayerIndex is the layer index
type LayerIndex int64

const (
	// ZIndexTop is to set the ZIndex at the top
	ZIndexTop = int64(math.MinInt64)
	// ZIndexBottom is to set the ZIndex at the bottom
	ZIndexBottom = int64(math.MinInt64 + 1)
)

type DrawLayer struct {
	Layer  LayerIndex
	ZIndex int64

	prevLayer LayerIndex
	prevIndex int64
}

//go:generate ecsgen -n DrawLayer -p graphics -o drawlayer_component.go --component-tpl --vars "UUID=2D35C735-7275-4195-A61F-F559F8346D46"

type drawLayerDrawer struct {
	index    LayerIndex
	slcindex int
	items    *core.EntitySortedList
}

type drawLayerItemCache struct {
	ZIndex    int64
	Entity    ecs.Entity
	Drawable  Drawable
	Transform *components.Transform
}

func (c *drawLayerItemCache) Less(v interface{}) bool {
	return c.ZIndex < v.(*drawLayerItemCache).ZIndex
}

func (c *drawLayerItemCache) Destroy() {
	c.Drawable = nil
}

type drawLayerDrawers struct {
	l     sync.RWMutex
	slice []*drawLayerDrawer
	m     map[LayerIndex]*drawLayerDrawer
}

func (d *drawLayerDrawers) UpsertLayer(index LayerIndex) *core.EntitySortedList {
	d.l.RLock()
	x := d.m[index]
	d.l.RUnlock()
	if x != nil {
		return x.items
	}
	d.l.Lock()
	defer d.l.Unlock()
	x = &drawLayerDrawer{
		index:    index,
		items:    core.NewEntitySortedList(1024),
		slcindex: len(d.slice),
	}
	d.m[index] = x
	d.slice = append(d.slice, x)
	sort.Sort(drawLayerDrawersL(d.slice))
	return x.items
}

// LayerTuple is returned when fetching all layers
type LayerTuple struct {
	Index LayerIndex
	Items *core.EntitySortedList
}

// All layers
func (d *drawLayerDrawers) All() []LayerTuple {
	d.l.RLock()
	defer d.l.RUnlock()
	cl := make([]LayerTuple, 0, len(d.slice))
	for _, v := range d.slice {
		cl = append(cl, LayerTuple{
			Index: v.index,
			Items: v.items,
		})
	}
	return cl
}

type drawLayerDrawersL []*drawLayerDrawer

func (a drawLayerDrawersL) Len() int { return len(a) }
func (a drawLayerDrawersL) Swap(i, j int) {
	a[i].slcindex, a[j].slcindex = j, i
	a[i], a[j] = a[j], a[i]
}
func (a drawLayerDrawersL) Less(i, j int) bool { return a[i].index < a[j].index }
