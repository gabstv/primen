package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

//FIXME: review (2)

// DrawableObj is the interface of a custom drawable component to use inetead
// of Drawable.Draw()
type DrawableObj interface {
	Draw(ctx DrawCtx, o *Drawable)
}

// Drawable is the component that controls drawing to the ebiten screen
type Drawable struct {
	concatm   ebiten.GeoM
	concatset bool
	drawer    DrawableObj
}

// Draw uses drawer or defaultDraw
func (d *Drawable) Draw(ctx DrawCtx) {
	if d.drawer != nil {
		d.drawer.Draw(ctx, d)
	}
}

func (d *Drawable) SetConcatM(m ebiten.GeoM) {
	d.concatm = m
	d.concatset = true
}

// G returns the concatm (if set by a *Transform) or a pre calculated matrix using
// sx, sy (scale)
// r (radians)
// x, y position
func (d *Drawable) G(sx, sy, r, x, y float64) ebiten.GeoM {
	if d.concatset {
		return d.concatm
	}
	return GeoM().Scale(sx, sy).Rotate(r).Translate(x, y).MV()
}

//go:generate ecsgen -n Drawable -p core -o drawable_component.go --component-tpl --vars "UUID=E3086C37-F0F5-4BFD-8FEE-F9C451B1E57E" --vars "OnWillResize=c.willresize()"

func (c *DrawableComponent) willresize() {
	for i := range c.data {
		c.data[i].Data.drawer = nil
	}
}

// ╔═╗╔═╗╦  ╔═╗  ╔═╗╦ ╦╔═╗
// ╚═╗║ ║║  ║ ║  ╚═╗╚╦╝╚═╗
// ╚═╝╚═╝╩═╝╚═╝  ╚═╝ ╩ ╚═╝

//go:generate ecsgen -n SoloDrawable -p core -o drawable_solosystem.go --system-tpl --vars "Priority=0" --vars "UUID=6389F54D-76C9-49FC-B3E3-1C73B334EBB6" --components "Drawable"

var matchSoloDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	if f.Contains(GetDrawLayerComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchSoloDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority is noop as of now
func (s *SoloDrawableSystem) DrawPriority(ctx DrawCtx) {}

// Draw all solo drawables ordered by entity ID
func (s *SoloDrawableSystem) Draw(ctx DrawCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.Draw(ctx)
	}
}

// UpdatePriority sets v.Drawable.concatset to false on all Drawables
func (s *SoloDrawableSystem) UpdatePriority(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.concatset = false
	}
}

// Update is noop as of now
func (s *SoloDrawableSystem) Update(ctx UpdateCtx) {}

//  ___   ___    __    _       _      __    _     ____  ___       __   _     __
// | | \ | |_)  / /\  \ \    /| |    / /\  \ \_/ | |_  | |_)     ( (` \ \_/ ( (`
// |_|_/ |_| \ /_/--\  \_\/\/ |_|__ /_/--\  |_|  |_|__ |_| \     _)_)  |_|  _)_)

//go:generate ecsgen -n DrawLayerDrawable -p core -o drawable_layersystem.go --system-tpl --vars "Priority=-10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "Setup=s.setupVars()" --vars "UUID=CBBC8DB4-4866-413E-A7A9-250A3C9ECDDC" --vars "OnWillResize=s.beforeCompResize()" --vars "OnResize=s.afterCompResize()" --components "Drawable" --components "DrawLayer" --members "layers=*drawLayerDrawers"

var matchDrawLayerDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawLayerComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchDrawLayerDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawLayerComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority is noop as of now
func (s *DrawLayerDrawableSystem) DrawPriority(ctx DrawCtx) {}

// Draw draws a drawable by its layer and zindex order
func (s *DrawLayerDrawableSystem) Draw(ctx DrawCtx) {
	for _, l := range s.layers.All() {
		l.Items.Each(func(key ecs.Entity, value SLVal) bool {
			cache := value.(*drawLayerItemCache)
			cache.Drawable.Draw(ctx)
			return true
		})
	}
}

// UpdatePriority updates layer changes
func (s *DrawLayerDrawableSystem) UpdatePriority(ctx UpdateCtx) {
	for _, v := range s.V().Matches() {
		// reset concat matrix:
		v.Drawable.concatset = false

		if v.DrawLayer.Layer != v.DrawLayer.prevLayer {
			// switch layers
			if l := s.layers.UpsertLayer(v.DrawLayer.prevLayer); l != nil {
				if x, ok := l.Get(v.Entity); ok {
					x.Destroy()
					l.Delete(v.Entity)
				}
			}
			v.DrawLayer.prevLayer = v.DrawLayer.Layer
			//
			l := s.layers.UpsertLayer(v.DrawLayer.Layer)
			// update index history since the layer changed anyway
			s.resolveIndex(v)
			v.DrawLayer.prevIndex = v.DrawLayer.ZIndex
			// TODO: check for leaks on AddOrUpdate (*Drawable might leak)
			l.AddOrUpdate(v.Entity, &drawLayerItemCache{
				ZIndex:   v.DrawLayer.ZIndex,
				Entity:   v.Entity,
				Drawable: v.Drawable,
			})
		} else if v.DrawLayer.ZIndex != v.DrawLayer.prevIndex {
			s.resolveIndex(v)
			v.DrawLayer.prevIndex = v.DrawLayer.ZIndex
			if l := s.layers.UpsertLayer(v.DrawLayer.Layer); l != nil {
				l.AddOrUpdate(v.Entity, &drawLayerItemCache{
					ZIndex:   v.DrawLayer.ZIndex,
					Drawable: v.Drawable,
				})
			}
		}
	}
}

// Update is noop as of now
func (s *DrawLayerDrawableSystem) Update(ctx UpdateCtx) {}

func (s *DrawLayerDrawableSystem) resolveIndex(v VIDrawLayerDrawableSystem) {
	if v.DrawLayer.ZIndex == ZIndexTop {
		l := s.layers.UpsertLayer(v.DrawLayer.Layer)
		if lv := l.LastValue(); lv != nil {
			v.DrawLayer.ZIndex = lv.(*drawLayerItemCache).ZIndex + 1
		} else {
			v.DrawLayer.ZIndex = 0
		}
	} else if v.DrawLayer.ZIndex == ZIndexBottom {
		l := s.layers.UpsertLayer(v.DrawLayer.Layer)
		if lv := l.FirstValue(); lv != nil {
			v.DrawLayer.ZIndex = lv.(*drawLayerItemCache).ZIndex - 1
		} else {
			v.DrawLayer.ZIndex = 0
		}
	}
}

func (s *DrawLayerDrawableSystem) onEntityAdded(e ecs.Entity) {
	dl := GetDrawLayerComponentData(s.world, e)
	d := GetDrawableComponentData(s.world, e)
	dl.prevLayer = dl.Layer
	l := s.layers.UpsertLayer(dl.Layer)
	if dl.ZIndex == ZIndexBottom {
		if lv := l.FirstValue(); lv != nil {
			dl.ZIndex = lv.(*drawLayerItemCache).ZIndex - 1
		} else {
			dl.ZIndex = 0
		}
	}
	if dl.ZIndex == ZIndexTop {
		if lv := l.LastValue(); lv != nil {
			dl.ZIndex = lv.(*drawLayerItemCache).ZIndex + 1
		} else {
			dl.ZIndex = 0
		}
	}
	dl.prevIndex = dl.ZIndex
	_ = l.AddOrUpdate(e, &drawLayerItemCache{
		ZIndex:   dl.ZIndex,
		Drawable: d,
	})
}

func (s *DrawLayerDrawableSystem) onEntityRemoved(e ecs.Entity) {
	for _, l := range s.layers.All() {
		if x, ok := l.Items.Get(e); ok {
			x.Destroy()
			l.Items.Delete(e)
		}
	}
}

func (s *DrawLayerDrawableSystem) setupVars() {
	s.layers = &drawLayerDrawers{
		slice: make([]*drawLayerDrawer, 0, 16),
		m:     make(map[LayerIndex]*drawLayerDrawer),
	}
}

func (s *DrawLayerDrawableSystem) beforeCompResize() {
	for _, l := range s.layers.All() {
		l.Items.Each(func(key ecs.Entity, value SLVal) bool {
			value.(*drawLayerItemCache).Drawable = nil
			return true
		})
	}
}

func (s *DrawLayerDrawableSystem) afterCompResize() {
	for _, l := range s.layers.All() {
		l.Items.Each(func(key ecs.Entity, value SLVal) bool {
			value.(*drawLayerItemCache).Drawable = GetDrawableComponentData(s.world, key)
			return true
		})
	}
}
