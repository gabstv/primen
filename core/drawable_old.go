package core

// import (
// 	"github.com/gabstv/ecs"
// 	"github.com/hajimehoshi/ebiten"
// )

// const (
// 	// DefaultImageOptions key passed to the default world (&ebiten.DrawImageOptions{})
// 	DefaultImageOptions string = "default_image_options"
// )

// // Drawable is the basis of all ebiten drawable items (sprites, texts, shapes)
// type Drawable interface {
// 	Update(ctx Context)
// 	Draw(m DrawManager)
// 	Destroy()
// 	IsDisabled() bool
// 	Size() (w, h float64)
// 	SetTransformMatrix(m GeoMatrix)
// 	ClearTransformMatrix()
// 	SetOffset(x, y float64)
// }

// // DrawableImager is a Drawable with Get and Set image funcs
// type DrawableImager interface {
// 	Drawable
// 	GetImage() *ebiten.Image
// 	SetImage(img *ebiten.Image)
// }

// const (
// 	// SNSoloDrawable is the system name of a drawable without as DrawLayer component
// 	SNSoloDrawable = "primen.DrawableSystem"
// 	// SNDrawLayerDrawable is the system name of a drawable with as DrawLayer component
// 	SNDrawLayerDrawable = "primen.DrawLayerDrawableSystem"
// 	// CNDrawable is the component name of a drawable
// 	CNDrawable = "primen.DrawableComponent"
// )

// // ███████╗ ██████╗ ██╗      ██████╗
// // ██╔════╝██╔═══██╗██║     ██╔═══██╗
// // ███████╗██║   ██║██║     ██║   ██║
// // ╚════██║██║   ██║██║     ██║   ██║
// // ███████║╚██████╔╝███████╗╚██████╔╝
// // ╚══════╝ ╚═════╝ ╚══════╝ ╚═════╝

// // SoloDrawableComponentSystem handles drawable items (without a draw layer)
// type SoloDrawableComponentSystem struct {
// 	BaseComponentSystem
// }

// // SystemName returns the system name
// func (cs *SoloDrawableComponentSystem) SystemName() string { return SNSoloDrawable }

// // SystemPriority returns the system priority
// func (cs *SoloDrawableComponentSystem) SystemPriority() int { return -10 }

// // SystemInit returns the system init
// func (cs *SoloDrawableComponentSystem) SystemInit() SystemInitFn {
// 	return func(w *ecs.World, sys *ecs.System) {
// 		if w.Get(DefaultImageOptions) == nil {
// 			opt := &ebiten.DrawImageOptions{}
// 			w.Set(DefaultImageOptions, opt)
// 		}
// 	}
// }

// // SystemExec returns the system exec fn
// func (cs *SoloDrawableComponentSystem) SystemExec() SystemExecFn {
// 	return soloDrawableComponentSystemExec
// }

// // Components returns the component signature(s)
// func (cs *SoloDrawableComponentSystem) Components(w *ecs.World) []*ecs.Component {
// 	return []*ecs.Component{
// 		drawableComponentDef(w),
// 	}
// }

// // ExcludeComponents returns the components that must not be present in this system
// func (cs *SoloDrawableComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
// 	return []*ecs.Component{
// 		drawLayerComponentDef(w),
// 	}
// }

// func drawableComponentDef(w *ecs.World) *ecs.Component {
// 	return UpsertComponent(w, ecs.NewComponentInput{
// 		Name: CNDrawable,
// 		ValidateDataFn: func(data interface{}) bool {
// 			_, ok := data.(Drawable)
// 			return ok
// 		},
// 		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
// 			sd := data.(Drawable)
// 			sd.Destroy()
// 		},
// 	})
// }

// // SystemTags -> draw
// func (cs *SoloDrawableComponentSystem) SystemTags() []string {
// 	return []string{"draw"}
// }

// // soloDrawableComponentSystemExec is the main function of the SoloDrawableComponentSystem
// func soloDrawableComponentSystemExec(ctx Context) {
// 	renderer := ctx.Renderer()
// 	v := ctx.System().View()
// 	matches := v.Matches()
// 	comp := ctx.World().Component(CNDrawable)
// 	for _, m := range matches {
// 		drawable := m.Components[comp].(Drawable)
// 		drawable.Update(ctx)
// 		if drawable.IsDisabled() {
// 			continue
// 		}
// 		drawable.Draw(renderer)
// 	}
// }

// // ██████╗ ██████╗  █████╗ ██╗    ██╗██╗      █████╗ ██╗   ██╗███████╗██████╗
// // ██╔══██╗██╔══██╗██╔══██╗██║    ██║██║     ██╔══██╗╚██╗ ██╔╝██╔════╝██╔══██╗
// // ██║  ██║██████╔╝███████║██║ █╗ ██║██║     ███████║ ╚████╔╝ █████╗  ██████╔╝
// // ██║  ██║██╔══██╗██╔══██║██║███╗██║██║     ██╔══██║  ╚██╔╝  ██╔══╝  ██╔══██╗
// // ██████╔╝██║  ██║██║  ██║╚███╔███╔╝███████╗██║  ██║   ██║   ███████╗██║  ██║
// // ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝

// // DrawLayerDrawableComponentSystem handles drawables with a draw layer
// type DrawLayerDrawableComponentSystem struct {
// 	BaseComponentSystem
// }

// // SystemName returns the system name
// func (cs *DrawLayerDrawableComponentSystem) SystemName() string {
// 	return SNDrawLayerDrawable
// }

// // SystemPriority returns the system priority
// func (cs *DrawLayerDrawableComponentSystem) SystemPriority() int {
// 	return -9
// }

// // SystemInit returns the system init
// func (cs *DrawLayerDrawableComponentSystem) SystemInit() SystemInitFn {
// 	return func(w *ecs.World, sys *ecs.System) {
// 		sys.View().SetOnEntityAdded(func(e ecs.Entity, w *ecs.World) {
// 			//TODO: checks?
// 		})
// 		sys.View().SetOnEntityRemoved(func(e ecs.Entity, w *ecs.World) {
// 			//TODO: checks?
// 		})
// 	}
// }

// func (cs *DrawLayerDrawableComponentSystem) SystemTags() []string {
// 	return []string{
// 		"draw",
// 	}
// }

// // SystemExec returns the system exec fn
// func (cs *DrawLayerDrawableComponentSystem) SystemExec() SystemExecFn {
// 	return drawLayerDrawableSystemExec
// }

// // Components returns the component signature(s)
// func (cs *DrawLayerDrawableComponentSystem) Components(w *ecs.World) []*ecs.Component {
// 	return []*ecs.Component{
// 		drawLayerComponentDef(w),
// 		drawableComponentDef(w),
// 	}
// }

// // drawLayerDrawableSystemExec is the main function of the DrawLayerSystem
// func drawLayerDrawableSystemExec(ctx Context) {
// 	world := ctx.World()
// 	renderer := ctx.Renderer()
// 	layers := world.System(SNDrawLayer).Get("layers").(*drawLayerDrawers).All()
// 	dwgetter := world.Component(CNDrawable)
// 	for _, layer := range layers {
// 		layer.Items.Each(func(key ecs.Entity, value SLVal) bool {
// 			cache := value.(*drawLayerItemCache)
// 			if cache.Drawable == nil {
// 				cache.Drawable = dwgetter.Data(key).(Drawable)
// 			}
// 			cache.Drawable.Update(ctx)
// 			if cache.Drawable.IsDisabled() {
// 				return true
// 			}
// 			cache.Drawable.Draw(renderer)
// 			return true
// 		})
// 	}
// }

// // ██╗███╗   ██╗██╗████████╗
// // ██║████╗  ██║██║╚══██╔══╝
// // ██║██╔██╗ ██║██║   ██║
// // ██║██║╚██╗██║██║   ██║
// // ██║██║ ╚████║██║   ██║
// // ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝

// func init() {
// 	RegisterComponentSystem(&SoloDrawableComponentSystem{})
// 	RegisterComponentSystem(&DrawLayerDrawableComponentSystem{})
// }
