package graphics

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/hajimehoshi/ebiten"
)

// Sprite is the data of a sprite component.
type Sprite struct {
	originX  float64 // X origin (0 = left; 0.5 = center; 1 = right)
	originY  float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	offsetX  float64 // offset origin X (in pixels)
	offsetY  float64 // offset origin Y (in pixels)
	image    *ebiten.Image
	disabled bool

	// is recalculated if image is set:

	imageWidth  float64 // last calculated image width
	imageHeight float64 // last calculated image height
	opt         ebiten.DrawImageOptions
}

func getImageSize(img *ebiten.Image) (w, h float64) {
	if img == nil {
		return 0, 0
	}
	iw, ih := img.Size()
	return float64(iw), float64(ih)
}

// NewSprite creates a new sprite (component data)
func NewSprite(x, y float64, quad *ebiten.Image) Sprite {
	iw, ih := getImageSize(quad)
	return Sprite{
		image:       quad,
		imageWidth:  iw,
		imageHeight: ih,
		opt:         ebiten.DrawImageOptions{},
	}
}

// public setters/getters

func (s *Sprite) SetEnabled(enabled bool) *Sprite {
	s.disabled = !enabled
	return s
}

func (s *Sprite) Origin() (ox, oy float64) {
	return s.originX, s.originY
}

func (s *Sprite) SetOrigin(ox, oy float64) *Sprite {
	s.originX, s.originY = ox, oy
	return s
}

func (s *Sprite) SetOffset(x, y float64) *Sprite {
	s.offsetX, s.offsetY = x, y
	return s
}

func (s *Sprite) ResetColorMatrix() {
	s.opt.ColorM.Reset()
}

func (s *Sprite) RotateHue(theta float64) {
	s.opt.ColorM.RotateHue(theta)
}

func (s *Sprite) SetCompositeMode(mode ebiten.CompositeMode) {
	s.opt.CompositeMode = mode
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) SetImage(img *ebiten.Image) *Sprite {
	s.image = img
	s.imageWidth, s.imageHeight = getImageSize(img)
	return s
}

func (s *Sprite) Update(ctx core.UpdateCtx, t *components.Transform) {}

func (s *Sprite) Draw(ctx core.DrawCtx, t *components.Transform) {
	if s.disabled {
		return
	}
	g := t.GeoM()
	o := &s.opt
	o.GeoM.Reset()
	o.GeoM.Translate(core.ApplyOrigin(s.imageWidth, s.originX)+s.offsetX, core.ApplyOrigin(s.imageHeight, s.originY)+s.offsetY)
	o.GeoM.Concat(g)

	//TODO: reimplement colormode and composite mode
	ctx.Renderer().DrawImageRaw(s.image, o)

	if debug.Draw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+s.imageWidth, y0
		x2, y2 := x1, y1+s.imageHeight
		x3, y3 := x2-s.imageWidth, y2
		screen := ctx.Renderer().Screen()
		debug.LineM(screen, o.GeoM, x0, y0, x1, y1, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x1, y1, x2, y2, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x2, y2, x3, y3, debug.BoundsColor)
		debug.LineM(screen, o.GeoM, x3, y3, x0, y0, debug.BoundsColor)
		debug.LineM(screen, g, -4, 0, 4, 0, debug.PivotColor)
		debug.LineM(screen, g, 0, -4, 0, 4, debug.PivotColor)
	}
}

//go:generate ecsgen -n Sprite -p graphics -o sprite_component.go --component-tpl --vars "UUID=80C95DEC-DBBF-4529-BD27-739A69055BA0" --vars "Setup=c.onCompSetup()"

func (c *SpriteComponent) onCompSetup() {
	RegisterDrawableComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) Drawable {
		return GetSpriteComponentData(w, e)
	})
}
