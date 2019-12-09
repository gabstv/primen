package troupe

import (
	"image"

	"github.com/gabstv/troupe/pkg/troupe/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	spriteComponentName = "troupe.Sprite"
)

const (
	// SpritePriority - default -10
	SpritePriority int = -10
	// DefaultImageOptions key passed to the default world (&ebiten.DrawImageOptions{})
	DefaultImageOptions string = "default_image_options"
)

func init() {
	DefaultComp(func(e *Engine, w *ecs.World) {
		SpriteComponent(w)
	})
	DefaultSys(func(e *Engine, w *ecs.World) {
		SpriteSystem(w)
	})
	println("graphicsinit end")
}

// Sprite is the data of a sprite component.
type Sprite struct {
	X      float64
	Y      float64
	Angle  float64
	ScaleX float64
	ScaleY float64
	// Bounds for drawing subimage
	Bounds image.Rectangle

	Options *ebiten.DrawImageOptions
	Image   *ebiten.Image
	// lastImage exists to keep track of the public Image field, if it
	// changes, the imageWidth and ImageHeight needs to be recalculated.
	lastImage   *ebiten.Image
	imageWidth  float64
	imageHeight float64

	lastBounds   image.Rectangle
	lastSubImage *ebiten.Image
}

// SpriteComponent will get the registered sprite component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func SpriteComponent(w *ecs.World) *ecs.Component {
	c := w.Component(spriteComponentName)
	if c == nil {
		var err error
		c, err = w.NewComponent(ecs.NewComponentInput{
			Name: spriteComponentName,
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*Sprite)
				return ok
			},
			DestructorFn: func(_ ecs.Worlder, entity ecs.Entity, data interface{}) {
				sd := data.(*Sprite)
				sd.Options = nil
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return c
}

// SpriteSystem creates the sprite system
func SpriteSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(SpritePriority, SpriteSystemExec, w.Component(spriteComponentName))
	if w.Get(DefaultImageOptions) == nil {
		opt := &ebiten.DrawImageOptions{}
		w.Set(DefaultImageOptions, opt)
	}
	sys.AddTag(WorldTagDraw)
	return sys
}

// SpriteSystemExec is the main function of the SpriteSystem
func SpriteSystemExec(ctx ecs.Context, screen *ebiten.Image) {
	// dt float64, v *ecs.View, s *ecs.System
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	spritecomp := world.Component(spriteComponentName)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	for _, m := range matches {
		sprite := m.Components[spritecomp].(*Sprite)
		opt := sprite.Options
		if opt == nil {
			opt = defaultopts
		}
		if sprite.lastImage != sprite.Image {
			w, h := sprite.Image.Size()
			sprite.imageWidth = float64(w)
			sprite.imageHeight = float64(h)
			sprite.lastImage = sprite.Image
			// redo subimage
			sprite.lastBounds = image.Rect(0, 0, 0, 0)
		}
		if sprite.lastBounds != sprite.Bounds {
			sprite.lastBounds = sprite.Bounds
			sprite.lastSubImage = sprite.Image.SubImage(sprite.lastBounds).(*ebiten.Image)
			w, h := sprite.lastSubImage.Size()
			sprite.imageWidth = float64(w)
			sprite.imageHeight = float64(h)
		}
		opt.GeoM.Reset()
		opt.GeoM.Translate(-sprite.imageWidth/2, -sprite.imageHeight/2)
		opt.GeoM.Scale(sprite.ScaleX, sprite.ScaleY)
		opt.GeoM.Rotate(sprite.Angle)
		opt.GeoM.Translate(sprite.imageWidth/2, sprite.imageHeight/2)
		opt.GeoM.Translate(sprite.X, sprite.Y)
		if sprite.lastSubImage != nil {
			screen.DrawImage(sprite.lastSubImage, opt)
		} else {
			screen.DrawImage(sprite.Image, opt)
		}
	}
}
