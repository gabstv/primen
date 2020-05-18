package tau

import (
	"image"

	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	// DefaultImageOptions key passed to the default world (&ebiten.DrawImageOptions{})
	DefaultImageOptions string = "default_image_options"
)

// Drawable is the basis of all ebiten drawable items (sprites, texts, shapes)
type Drawable interface {
	Update(ctx Context)
	Draw(screen *ebiten.Image, opt *ebiten.DrawImageOptions)
	Destroy()
	DrawImageOptions() *ebiten.DrawImageOptions
	IsDisabled() bool
	Size() (w, h float64)
	SetTransformMatrix(m ebiten.GeoM)
	SetBounds(b image.Rectangle)
}

const (
	// SNSoloDrawable is the system name of a drawable without as DrawLayer component
	SNSoloDrawable = "tau.DrawableSystem"
	// SNDrawLayerDrawable is the system name of a drawable with as DrawLayer component
	SNDrawLayerDrawable = "tau.DrawLayerDrawableSystem"
	// CNDrawable is the component name of a drawable
	CNDrawable = "tau.DrawableComponent"
)

// ███████╗ ██████╗ ██╗      ██████╗
// ██╔════╝██╔═══██╗██║     ██╔═══██╗
// ███████╗██║   ██║██║     ██║   ██║
// ╚════██║██║   ██║██║     ██║   ██║
// ███████║╚██████╔╝███████╗╚██████╔╝
// ╚══════╝ ╚═════╝ ╚══════╝ ╚═════╝

// SoloDrawableComponentSystem handles drawable items (without a draw layer)
type SoloDrawableComponentSystem struct {
	BaseComponentSystem
}

// SystemName returns the system name
func (cs *SoloDrawableComponentSystem) SystemName() string { return SNSoloDrawable }

// SystemPriority returns the system priority
func (cs *SoloDrawableComponentSystem) SystemPriority() int { return -10 }

// SystemInit returns the system init
func (cs *SoloDrawableComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		if w.Get(DefaultImageOptions) == nil {
			opt := &ebiten.DrawImageOptions{}
			w.Set(DefaultImageOptions, opt)
		}
	}
}

// SystemExec returns the system exec fn
func (cs *SoloDrawableComponentSystem) SystemExec() SystemExecFn {
	return soloDrawableComponentSystemExec
}

// Components returns the component signature(s)
func (cs *SoloDrawableComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawableComponentDef(w),
	}
}

// ExcludeComponents returns the components that must not be present in this system
func (cs *SoloDrawableComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawLayerComponentDef(w),
	}
}

func drawableComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNDrawable,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(Drawable)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			sd := data.(Drawable)
			sd.Destroy()
		},
	})
}

// SystemTags -> draw
func (cs *SoloDrawableComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

// soloDrawableComponentSystemExec is the main function of the SoloDrawableComponentSystem
func soloDrawableComponentSystemExec(ctx Context) {
	screen := ctx.Screen()
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	comp := ctx.World().Component(CNDrawable)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	for _, m := range matches {
		drawable := m.Components[comp].(Drawable)
		drawable.Update(ctx)
		if drawable.IsDisabled() {
			continue
		}
		opt := drawable.DrawImageOptions()
		if opt == nil {
			opt = defaultopts
		}
		drawable.Draw(screen, opt)
	}
}

// ██████╗ ██████╗  █████╗ ██╗    ██╗██╗      █████╗ ██╗   ██╗███████╗██████╗
// ██╔══██╗██╔══██╗██╔══██╗██║    ██║██║     ██╔══██╗╚██╗ ██╔╝██╔════╝██╔══██╗
// ██║  ██║██████╔╝███████║██║ █╗ ██║██║     ███████║ ╚████╔╝ █████╗  ██████╔╝
// ██║  ██║██╔══██╗██╔══██║██║███╗██║██║     ██╔══██║  ╚██╔╝  ██╔══╝  ██╔══██╗
// ██████╔╝██║  ██║██║  ██║╚███╔███╔╝███████╗██║  ██║   ██║   ███████╗██║  ██║
// ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝

// DrawLayerDrawableComponentSystem handles drawables with a draw layer
type DrawLayerDrawableComponentSystem struct {
	BaseComponentSystem
}

// SystemName returns the system name
func (cs *DrawLayerDrawableComponentSystem) SystemName() string {
	return SNDrawLayerDrawable
}

// SystemPriority returns the system priority
func (cs *DrawLayerDrawableComponentSystem) SystemPriority() int {
	return -9
}

// SystemInit returns the system init
func (cs *DrawLayerDrawableComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		sys.View().SetOnEntityAdded(func(e ecs.Entity, w *ecs.World) {
			//TODO: checks?
		})
		sys.View().SetOnEntityRemoved(func(e ecs.Entity, w *ecs.World) {
			//TODO: checks?
		})
	}
}

// SystemExec returns the system exec fn
func (cs *DrawLayerDrawableComponentSystem) SystemExec() SystemExecFn {
	return drawLayerDrawableSystemExec
}

// Components returns the component signature(s)
func (cs *DrawLayerDrawableComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawLayerComponentDef(w),
		drawableComponentDef(w),
	}
}

// drawLayerDrawableSystemExec is the main function of the DrawLayerSystem
func drawLayerDrawableSystemExec(ctx Context) {
	world := ctx.World()
	screen := ctx.Screen()
	layers := world.System(SNDrawLayer).Get("layers").(*drawLayerDrawers).All()
	dwgetter := world.Component(CNDrawable)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	for _, layer := range layers {
		layer.Items.Each(func(key ecs.Entity, value SLVal) bool {
			cache := value.(*drawLayerItemCache)
			if cache.Drawable == nil {
				cache.Drawable = dwgetter.Data(key).(Drawable)
			}
			cache.Drawable.Update(ctx)
			if cache.Drawable.IsDisabled() {
				return true
			}
			opt := cache.Drawable.DrawImageOptions()
			if opt == nil {
				opt = defaultopts
			}
			cache.Drawable.Draw(screen, opt)
			return true
		})
	}
}

// ██╗███╗   ██╗██╗████████╗
// ██║████╗  ██║██║╚══██╔══╝
// ██║██╔██╗ ██║██║   ██║
// ██║██║╚██╗██║██║   ██║
// ██║██║ ╚████║██║   ██║
// ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝

func init() {
	RegisterComponentSystem(&SoloDrawableComponentSystem{})
	RegisterComponentSystem(&DrawLayerDrawableComponentSystem{})
}
