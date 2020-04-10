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

const (
	SNSprite = "tau.SpriteSystem"
	CNSprite = "tau.SpriteComponent"
)

var (
	SpriteCS *SpriteComponentSystem = new(SpriteComponentSystem)
)

type SpriteComponentSystem struct {
	BaseComponentSystem
}

func (cs *SpriteComponentSystem) SystemName() string {
	return SNSprite
}

func (cs *SpriteComponentSystem) SystemPriority() int {
	return -10
}

func (cs *SpriteComponentSystem) SystemExec() SystemExecFn {
	return SpriteSystemExec
}

func (cs *SpriteComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		spriteComponentDef(w),
	}
}

func spriteComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNSprite,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Sprite)
			return ok
		},
		DestructorFn: func(_ ecs.WorldDicter, entity ecs.Entity, data interface{}) {
			sd := data.(*Sprite)
			sd.Options = nil
		},
	})
}

func (cs *SpriteComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

// func init() {
// 	DefaultComp(func(e *Engine, w *World) {
// 		SpriteComponent(w)
// 	})
// 	DefaultSys(func(e *Engine, w *World) {
// 		SpriteSystem(w)
// 	})
// 	println("graphicsinit end")
// }

// Sprite is the data of a sprite component.
type Sprite struct {
	X       float64
	Y       float64
	Angle   float64
	ScaleX  float64
	ScaleY  float64
	OriginX float64
	OriginY float64

	Bounds image.Rectangle // Bounds for drawing subimage

	Options *ebiten.DrawImageOptions
	Image   *ebiten.Image

	DrawDisabled bool // if true, the SpriteSystem will not draw this

	lastImage *ebiten.Image // lastImage exists to keep track of the public Image field, if it
	// changes, the imageWidth and ImageHeight needs to be recalculated.
	imageWidth  float64 // last calculated image width
	imageHeight float64 // last calculated image height

	lastBounds   image.Rectangle
	lastSubImage *ebiten.Image
}

// GetPrecomputedImage returns the last precomputed image
func (s *Sprite) GetPrecomputedImage() *ebiten.Image {
	if s.lastSubImage != nil {
		return s.lastSubImage
	}
	return s.lastImage
}

// GetPrecomputedImageDim returns the last precomputed image dimmensions
func (s *Sprite) GetPrecomputedImageDim() (width, height float64) {
	return s.imageWidth, s.imageHeight
}

// SpriteSystemExec is the main function of the SpriteSystem
func SpriteSystemExec(ctx Context) {
	screen := ctx.Screen()
	// dt float64, v *ecs.View, s *ecs.System
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	spritecomp := ctx.World().Component(CNSprite)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	hw, hh := 0.0, 0.0
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
		if sprite.DrawDisabled {
			continue
		}
		hw, hh = sprite.imageWidth/2, sprite.imageHeight/2
		opt.GeoM.Reset()
		opt.GeoM.Translate(-hw+sprite.OriginX*sprite.imageWidth*-1, -hh+sprite.OriginY*sprite.imageHeight*-1)
		opt.GeoM.Scale(sprite.ScaleX, sprite.ScaleY)
		opt.GeoM.Rotate(sprite.Angle)
		opt.GeoM.Translate(hw, hh)
		opt.GeoM.Translate(sprite.X, sprite.Y)
		if sprite.lastSubImage != nil {
			screen.DrawImage(sprite.lastSubImage, opt)
		} else {
			screen.DrawImage(sprite.Image, opt)
		}
	}
}
