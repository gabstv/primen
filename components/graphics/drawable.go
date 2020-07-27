package graphics

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
)

// Transform;*components.Transform;components.GetTransformComponentData(v.world, e)

// Drawable is the component that controls drawing to the ebiten screen
type Drawable interface {
	DrawMask() core.DrawMask
	SetDrawMask(mask core.DrawMask)
	Draw(ctx core.DrawCtx, t *components.Transform)
	Update(ctx core.UpdateCtx, t *components.Transform)
}

type GetDrawableFn func(w ecs.BaseWorld, e ecs.Entity) Drawable

func RegisterDrawableComponent(w ecs.BaseWorld, f ecs.Flag, fn GetDrawableFn) {
	gflags := w.FlagGroup("PrimenDrawables")
	gflags = gflags.Or(f)
	w.SetFlagGroup("PrimenDrawables", gflags)
	vi := w.LGet("PrimenDrawables")
	var vs map[uint8]GetDrawableFn
	if vi != nil {
		vs = vi.(map[uint8]GetDrawableFn)
	} else {
		vs = make(map[uint8]GetDrawableFn)
	}
	vs[f.Lowest()] = fn
	w.LSet("PrimenDrawables", vs)
}

func GetDrawable(w ecs.BaseWorld, e ecs.Entity) Drawable {
	eflag := w.CFlag(e)
	dflag := eflag.And(w.FlagGroup("PrimenDrawables"))
	vi := w.LGet("PrimenDrawables").(map[uint8]GetDrawableFn)
	getter := vi[dflag.Lowest()]
	return getter(w, e)
}

// ╔═╗╔═╗╦  ╔═╗  ╔═╗╦ ╦╔═╗
// ╚═╗║ ║║  ║ ║  ╚═╗╚╦╝╚═╗
// ╚═╝╚═╝╩═╝╚═╝  ╚═╝ ╩ ╚═╝

//go:generate ecsgen -n SoloDrawable -p graphics -o drawable_solosystem.go --system-tpl --vars "Priority=0" --vars "UUID=6389F54D-76C9-49FC-B3E3-1C73B334EBB6" --components "Drawable;Drawable;GetDrawable(v.world, e)" --components "Transform;*components.Transform;components.GetTransformComponentData(v.world, e)" --go-import "\"github.com/gabstv/primen/components\""

var matchSoloDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawLayerComponent(w).Flag()) {
		return false
	}
	if !f.Contains(components.GetTransformComponent(w).Flag()) {
		return false
	}
	if f.ContainsAny(w.FlagGroup("PrimenDrawables")) {
		return true
	}
	return false
}

var resizematchSoloDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.ContainsAny(w.FlagGroup("PrimenDrawables")) {
		return true
	}
	if f.Contains(components.GetTransformComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority is noop as of now
func (s *SoloDrawableSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw all solo drawables ordered by entity ID
func (s *SoloDrawableSystem) Draw(ctx core.DrawCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.Draw(ctx, v.Transform)
	}
}

// UpdatePriority noop
func (s *SoloDrawableSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update calls drawable.Update()
func (s *SoloDrawableSystem) Update(ctx core.UpdateCtx) {
	for _, v := range s.V().Matches() {
		v.Drawable.Update(ctx, v.Transform)
	}
}

//  ___   ___    __    _       _      __    _     ____  ___       __   _     __
// | | \ | |_)  / /\  \ \    /| |    / /\  \ \_/ | |_  | |_)     ( (` \ \_/ ( (`
// |_|_/ |_| \ /_/--\  \_\/\/ |_|__ /_/--\  |_|  |_|__ |_| \     _)_)  |_|  _)_)

//go:generate ecsgen -n DrawLayerDrawable -p graphics -o drawable_layersystem.go --system-tpl --vars "Priority=-10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "Setup=s.setupVars()" --vars "UUID=CBBC8DB4-4866-413E-A7A9-250A3C9ECDDC" --vars "OnWillResize=s.beforeCompResize()" --vars "OnResize=s.afterCompResize()" --components "Drawable;Drawable;GetDrawable(v.world, e)" --components "DrawLayer" --components "Transform;*components.Transform;components.GetTransformComponentData(v.world, e)" --go-import "\"github.com/gabstv/primen/components\"" --members "layers=*drawLayerDrawers"

var matchDrawLayerDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawLayerComponent(w).Flag()) {
		return false
	}
	if !f.Contains(components.GetTransformComponent(w).Flag()) {
		return false
	}
	if f.ContainsAny(w.FlagGroup("PrimenDrawables")) {
		return true
	}
	return false
}

var resizematchDrawLayerDrawableSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawLayerComponent(w).Flag()) {
		return true
	}
	if f.Contains(components.GetTransformComponent(w).Flag()) {
		return true
	}
	if f.ContainsAny(w.FlagGroup("PrimenDrawables")) {
		return true
	}
	return false
}

// DrawPriority is noop as of now
func (s *DrawLayerDrawableSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw draws a drawable by its layer and zindex order
func (s *DrawLayerDrawableSystem) Draw(ctx core.DrawCtx) {
	for _, l := range s.layers.All() {
		l.Items.Each(func(key ecs.Entity, value core.SLVal) bool {
			cache := value.(*drawLayerItemCache)
			cache.Drawable.Draw(ctx, cache.Transform)
			return true
		})
	}
}

// UpdatePriority updates layer changes
func (s *DrawLayerDrawableSystem) UpdatePriority(ctx core.UpdateCtx) {
	for _, v := range s.V().Matches() {
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
			// TODO: check for leaks on AddOrUpdate (Drawable might leak)
			l.AddOrUpdate(v.Entity, &drawLayerItemCache{
				ZIndex:    v.DrawLayer.ZIndex,
				Entity:    v.Entity,
				Drawable:  v.Drawable,
				Transform: v.Transform,
			})
		} else if v.DrawLayer.ZIndex != v.DrawLayer.prevIndex {
			s.resolveIndex(v)
			v.DrawLayer.prevIndex = v.DrawLayer.ZIndex
			if l := s.layers.UpsertLayer(v.DrawLayer.Layer); l != nil {
				l.AddOrUpdate(v.Entity, &drawLayerItemCache{
					ZIndex:    v.DrawLayer.ZIndex,
					Drawable:  v.Drawable,
					Transform: v.Transform,
				})
			}
		}
	}
}

// Update is noop as of now
func (s *DrawLayerDrawableSystem) Update(ctx core.UpdateCtx) {}

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
	d := GetDrawable(s.world, e)
	tr := components.GetTransformComponentData(s.world, e)
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
		ZIndex:    dl.ZIndex,
		Drawable:  d,
		Transform: tr,
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
		l.Items.Each(func(key ecs.Entity, value core.SLVal) bool {
			value.(*drawLayerItemCache).Drawable = nil
			value.(*drawLayerItemCache).Transform = nil
			return true
		})
	}
}

func (s *DrawLayerDrawableSystem) afterCompResize() {
	for _, l := range s.layers.All() {
		l.Items.Each(func(key ecs.Entity, value core.SLVal) bool {
			value.(*drawLayerItemCache).Drawable = GetDrawable(s.world, key)
			value.(*drawLayerItemCache).Transform = components.GetTransformComponentData(s.world, key)
			return true
		})
	}
}
