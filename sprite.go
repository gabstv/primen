package tau

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

/*
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

func (cs *SpriteComponentSystem) SystemInit() SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {
		if w.Get(DefaultImageOptions) == nil {
			opt := &ebiten.DrawImageOptions{}
			w.Set(DefaultImageOptions, opt)
		}
	}
}

func (cs *SpriteComponentSystem) SystemExec() SystemExecFn {
	return SpriteSystemExec
}

func (cs *SpriteComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		spriteComponentDef(w),
	}
}

func (cs *SpriteComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		drawLayerComponentDef(w),
	}
}

func spriteComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNSprite,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*Sprite)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			sd := data.(*Sprite)
			sd.Options = nil
		},
	})
}

func (cs *SpriteComponentSystem) SystemTags() []string {
	return []string{"draw"}
}*/

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

	imageBounds  image.Rectangle
	lastBounds   image.Rectangle
	lastSubImage *ebiten.Image
	//
	transformMatrix ebiten.GeoM
	customMatrix    bool
}

func (s *Sprite) Update(ctx Context) {
	if s.lastImage != s.Image {
		w, h := s.Image.Size()
		s.imageWidth = float64(w)
		s.imageHeight = float64(h)
		s.lastImage = s.Image
		// redo subimage
		s.lastBounds = image.Rect(0, 0, 0, 0)
		s.imageBounds = s.Image.Bounds()
	}
	if s.lastBounds != s.Bounds {
		s.lastBounds = s.Bounds
		if s.imageBounds.Min.Eq(s.Bounds.Min) && s.imageBounds.Max.Eq(s.Bounds.Max) {
			s.imageWidth = float64(s.Bounds.Dx())
			s.imageHeight = float64(s.Bounds.Dy())
			s.lastSubImage = nil
		} else {
			s.lastSubImage = s.Image.SubImage(s.lastBounds).(*ebiten.Image)
			w, h := s.lastSubImage.Size()
			s.imageWidth = float64(w)
			s.imageHeight = float64(h)
		}
	}
}

func (s *Sprite) Draw(screen *ebiten.Image, opt *ebiten.DrawImageOptions) {
	if s.DrawDisabled {
		return
	}
	prevGeo := opt.GeoM
	if s.customMatrix {
		opt.GeoM = s.transformMatrix
	} else {
		opt.GeoM.Scale(s.ScaleX, s.ScaleY)
		opt.GeoM.Rotate(s.Angle)
		opt.GeoM.Translate(s.X, s.Y)
	}
	xxg := &ebiten.GeoM{}
	xxg.Translate(applyOrigin(s.imageWidth, s.OriginX), applyOrigin(s.imageHeight, s.OriginY))
	xxg.Concat(opt.GeoM)
	opt.GeoM = *xxg
	//opt.GeoM.Translate(applyOrigin(s.imageWidth, s.OriginX), applyOrigin(s.imageHeight, s.OriginY))
	if s.lastSubImage != nil {
		screen.DrawImage(s.lastSubImage, opt)
	} else {
		screen.DrawImage(s.Image, opt)
	}
	if DebugDraw {
		x0, y0 := 0.0, 0.0 //applyOrigin(s.imageWidth, s.OriginX), applyOrigin(s.imageHeight, s.OriginY)
		x1, y1 := x0+s.imageWidth, y0
		x2, y2 := x1, y1+s.imageHeight
		x3, y3 := x2-s.imageWidth, y2
		debugLineM(screen, opt, x0, y0, x1, y1, debugBoundsColor)
		debugLineM(screen, opt, x1, y1, x2, y2, debugBoundsColor)
		debugLineM(screen, opt, x2, y2, x3, y3, debugBoundsColor)
		debugLineM(screen, opt, x3, y3, x0, y0, debugBoundsColor)
	}
	opt.GeoM = prevGeo
	s.customMatrix = false
}

func (s *Sprite) SetTransformMatrix(m ebiten.GeoM) {
	s.transformMatrix = m
	s.customMatrix = true
}

func (s *Sprite) Destroy() {
	s.Image = nil
	s.Options = nil
}

func (s *Sprite) DrawImageOptions() *ebiten.DrawImageOptions {
	return s.Options
}

func (s *Sprite) IsDisabled() bool {
	return s.DrawDisabled
}

func (s *Sprite) Size() (w, h float64) {
	return s.imageWidth, s.imageHeight
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

/*
func drawSprite(screen *ebiten.Image, spriteComp *ecs.Component, sprite *Sprite, opt *ebiten.DrawImageOptions) {
	if sprite.lastImage != sprite.Image {
		w, h := sprite.Image.Size()
		sprite.imageWidth = float64(w)
		sprite.imageHeight = float64(h)
		sprite.lastImage = sprite.Image
		// redo subimage
		sprite.lastBounds = image.Rect(0, 0, 0, 0)
		sprite.imageBounds = sprite.Image.Bounds()
	}
	if sprite.lastBounds != sprite.Bounds {
		sprite.lastBounds = sprite.Bounds
		if sprite.imageBounds.Min.Eq(sprite.Bounds.Min) && sprite.imageBounds.Max.Eq(sprite.Bounds.Max) {
			sprite.imageWidth = float64(sprite.Bounds.Dx())
			sprite.imageHeight = float64(sprite.Bounds.Dy())
			sprite.lastSubImage = nil
		} else {
			sprite.lastSubImage = sprite.Image.SubImage(sprite.lastBounds).(*ebiten.Image)
			w, h := sprite.lastSubImage.Size()
			sprite.imageWidth = float64(w)
			sprite.imageHeight = float64(h)
		}
	}
	if sprite.DrawDisabled {
		return
	}
	hw, hh := sprite.imageWidth/2, sprite.imageHeight/2
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

// SpriteSystemExec is the main function of the SpriteSystem
func SpriteSystemExec(ctx Context) {
	screen := ctx.Screen()
	// dt float64, v *ecs.View, s *ecs.System
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	spritecomp := ctx.World().Component(CNSprite)
	defaultopts := world.Get(DefaultImageOptions).(*ebiten.DrawImageOptions)
	for _, m := range matches {
		sprite := m.Components[spritecomp].(*Sprite)
		opt := sprite.Options
		if opt == nil {
			opt = defaultopts
		}
		drawSprite(screen, spritecomp, sprite, opt)
	}
}

func init() {
	RegisterComponentSystem(&SpriteComponentSystem{})
}
*/
