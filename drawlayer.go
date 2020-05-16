package tau

import (
	"math"
	"sort"
	"sync"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	// SNDrawLayer is the system name
	SNDrawLayer       = "tau.DrawLayerSystem"
	CNDrawLayer       = "tau.DrawLayerComponent"
	SNDrawLayerSprite = "tau.DrawLayerSpriteSystem"
)

type DrawLayer struct {
	Layer  LayerIndex
	ZIndex int64

	prevLayer LayerIndex
	prevIndex int64
}

const (
	// ZIndexLast is to set the ZIndex at the top
	ZIndexTop = int64(math.MinInt64)
	// ZIndexBottom is to set the ZIndex at the bottom
	ZIndexBottom = int64(math.MinInt64 + 1)
)

type LayerIndex int64

// LAYERS
const (
	Layer0  LayerIndex = 0
	Layer1  LayerIndex = 1
	Layer2  LayerIndex = 2
	Layer3  LayerIndex = 3
	Layer4  LayerIndex = 4
	Layer5  LayerIndex = 5
	Layer6  LayerIndex = 6
	Layer7  LayerIndex = 7
	Layer8  LayerIndex = 8
	Layer9  LayerIndex = 9
	Layer10 LayerIndex = 10
	Layer11 LayerIndex = 11
	Layer12 LayerIndex = 12
	Layer13 LayerIndex = 13
	Layer14 LayerIndex = 14
	Layer15 LayerIndex = 15
)

type drawLayerDrawer struct {
	index    LayerIndex
	slcindex int
	items    *EntitySortedList
}

type drawLayerItemCache struct {
	ZIndex int64
	Sprite *Sprite
}

func (c *drawLayerItemCache) Less(v interface{}) bool {
	return c.ZIndex < v.(*drawLayerItemCache).ZIndex
}

//
type drawLayerDrawers struct {
	l     sync.RWMutex
	slice []*drawLayerDrawer
	m     map[LayerIndex]*drawLayerDrawer
}

func (d *drawLayerDrawers) UpsertLayer(index LayerIndex) *EntitySortedList {
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
		items:    NewEntitySortedList(1024),
		slcindex: len(d.slice),
	}
	d.m[index] = x
	d.slice = append(d.slice, x)
	sort.Sort(drawLayerDrawersL(d.slice))
	return x.items
}

type LayerTuple struct {
	Index LayerIndex
	Items *EntitySortedList
}

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

func drawLayerComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNDrawLayer,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*DrawLayer)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			//sd := data.(*DrawLayer)
		},
	})
}

type DrawLayerComponentSystem struct {
	BaseComponentSystem
}

func (cs *DrawLayerComponentSystem) SystemName() string {
	return SNDrawLayer
}

func (cs *DrawLayerComponentSystem) SystemPriority() int {
	return 0
}

func (cs *DrawLayerComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		layers := &drawLayerDrawers{
			slice: make([]*drawLayerDrawer, 0, 16),
			m:     make(map[LayerIndex]*drawLayerDrawer),
		}
		sys.Set("layers", layers)
		sys.View().SetOnEntityAdded(func(e ecs.Entity, w *ecs.World) {
			// entity added to the system's view
			layers := sys.Get("layers").(*drawLayerDrawers)
			lcomp := w.Component(CNDrawLayer).Data(e).(*DrawLayer)
			lcomp.prevLayer = lcomp.Layer
			l := layers.UpsertLayer(lcomp.Layer)
			if lcomp.ZIndex == ZIndexBottom {
				if lv := l.FirstValue(); lv != nil {
					lcomp.ZIndex = lv.(*drawLayerItemCache).ZIndex - 1
				} else {
					lcomp.ZIndex = 0
				}
			}
			if lcomp.ZIndex == ZIndexTop {
				if lv := l.LastValue(); lv != nil {
					lcomp.ZIndex = lv.(*drawLayerItemCache).ZIndex + 1
				} else {
					lcomp.ZIndex = 0
				}
			}
			lcomp.prevIndex = lcomp.ZIndex
			_ = l.AddOrUpdate(e, &drawLayerItemCache{
				ZIndex: lcomp.ZIndex,
			})
		})
		sys.View().SetOnEntityRemoved(func(e ecs.Entity, w *ecs.World) {
			// entity removed from the system's view
			layers := sys.Get("layers").(*drawLayerDrawers)
			lcomp := w.Component(CNDrawLayer).Data(e).(*DrawLayer)
			if l := layers.UpsertLayer(lcomp.Layer); l != nil {
				_ = l.Delete(e)
			}
			if lcomp.prevLayer != lcomp.Layer {
				if l := layers.UpsertLayer(lcomp.prevLayer); l != nil {
					_ = l.Delete(e)
				}
			}
		})
	}
}

func (cs *DrawLayerComponentSystem) SystemExec() SystemExecFn {
	return DrawLayerSystemExec
}

func (cs *DrawLayerComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawLayerComponentDef(w),
	}
}

func (cs *DrawLayerComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

// DrawLayerSystemExec is the main function of the DrawLayerSystem
func DrawLayerSystemExec(ctx Context) {
	// dt float64, v *ecs.View, s *ecs.System
	s := ctx.System()
	v := s.View()
	layers := s.Get("layers").(*drawLayerDrawers)
	for _, match := range v.Matches() {
		dlc := match.Components[s.World().Component(CNDrawLayer)].(*DrawLayer)
		if dlc.Layer != dlc.prevLayer {
			// switch layers
			if l := layers.UpsertLayer(dlc.prevLayer); l != nil {
				_ = l.Delete(match.Entity)
			}
			dlc.prevLayer = dlc.Layer
			//
			l := layers.UpsertLayer(dlc.Layer)
			// update index history since the layer changed anyway
			dlc.prevIndex = dlc.ZIndex
			l.AddOrUpdate(match.Entity, &drawLayerItemCache{
				ZIndex: dlc.ZIndex,
			})
		} else if dlc.ZIndex != dlc.prevIndex {
			dlc.prevIndex = dlc.ZIndex
			if l := layers.UpsertLayer(dlc.Layer); l != nil {
				l.AddOrUpdate(match.Entity, &drawLayerItemCache{
					ZIndex: dlc.ZIndex,
				})
			}
		}
	}
}

//
//
//

type DrawLayerSpriteComponentSystem struct {
	BaseComponentSystem
}

func (cs *DrawLayerSpriteComponentSystem) SystemName() string {
	return SNDrawLayerSprite
}

func (cs *DrawLayerSpriteComponentSystem) SystemPriority() int {
	return -9
}

func (cs *DrawLayerSpriteComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		sys.View().SetOnEntityAdded(func(e ecs.Entity, w *ecs.World) {

		})
		sys.View().SetOnEntityRemoved(func(e ecs.Entity, w *ecs.World) {

		})
	}
}

func (cs *DrawLayerSpriteComponentSystem) SystemExec() SystemExecFn {
	return DrawLayerSpriteSystemExec
}

func (cs *DrawLayerSpriteComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawLayerComponentDef(w),
		spriteComponentDef(w),
	}
}

// DrawLayerSpriteSystemExec is the main function of the DrawLayerSystem
func DrawLayerSpriteSystemExec(ctx Context) {
	world := ctx.World()
	screen := ctx.Screen()
	layers := world.System(SNDrawLayer).Get("layers").(*drawLayerDrawers).All()
	spritec := world.Component(CNSprite)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	for _, layer := range layers {
		layer.Items.Each(func(key ecs.Entity, value SLVal) bool {
			cache := value.(*drawLayerItemCache)
			if cache.Sprite == nil {
				cache.Sprite = spritec.Data(key).(*Sprite)
			}
			opt := cache.Sprite.Options
			if opt == nil {
				opt = defaultopts
			}
			drawSprite(screen, spritec, cache.Sprite, opt)
			return true
		})
	}
}

func init() {
	RegisterComponentSystem(&DrawLayerComponentSystem{})
	RegisterComponentSystem(&DrawLayerSpriteComponentSystem{})
}
