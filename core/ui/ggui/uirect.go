package ggui

import (
	"image"
	"image/color"

	"github.com/gabstv/primen/core"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/internal/ggebiten"

	"github.com/hajimehoshi/ebiten"
)

type ColorMapFn func(x, y, width, height int) color.RGBA

type UIRect struct {
	filter      ebiten.Filter
	opt         ebiten.DrawImageOptions
	cache       *ebiten.Image
	notdirty    bool
	strokeColor color.RGBA
	stroke      int
	size        image.Point
	bgColor     color.RGBA
	bgColorFn   ColorMapFn
}

func (r *UIRect) render() {
	if r.cache != nil {
		r.cache.Dispose()
		r.cache = nil
	}
	if r.size.X == 0 || r.size.Y == 0 {
		return
	}
	gfx := ggebiten.NewGraphicsSoftLink(r.size.X+2+r.stroke*2, r.size.Y+2+r.stroke*2, r.filter)
	r.cache = gfx.Ebimage()
	//TODO: support gradients
	// https://github.com/fogleman/gg#gradients--patterns
	gfx.DrawRect(1+r.stroke, 1+r.stroke, r.size.X, r.size.Y, r.stroke, r.strokeColor, r.bgColor)
	gfx.Sync()
	gfx.Dispose()
}

func (r *UIRect) Draw(ctx core.DrawCtx, t *core.Transform) {
	if r.cache == nil {
		return
	}
	r.opt.GeoM.Reset()
	r.opt.GeoM.Translate(-1-float64(r.stroke), -1-float64(r.stroke))
	r.opt.GeoM.Concat(t.GeoM())
	ctx.Renderer().DrawImageRaw(r.cache, &r.opt)
}

func (r *UIRect) Update(ctx core.UpdateCtx, t *core.Transform) {
	if !r.notdirty {
		r.render()
		r.notdirty = true
	}
}

func (r *UIRect) SetSize(w, h float64) {
	r.size = image.Point{
		X: int(w),
		Y: int(h),
	}
}

//go:generate ecsgen -n UIRect -p core -o uirect_component.go --component-tpl --vars "UUID=5895E095-CECF-4AA0-A3A2-44460FDFC3FB" --vars "Setup=c.onCompSetup()"

func (c *UIRectComponent) onCompSetup() {
	core.RegisterDrawableComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) core.Drawable {
		return GetUIRectComponentData(w, e)
	})
	RegisterSizeableUIComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) SizeableUI {
		return GetUIRectComponentData(w, e)
	})
}