package primen

import (
	"image/color"

	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
)

var (
	transparentPixel *ebiten.Image
)

type Sprite struct {
	*WorldItem
	*TransformItem
	*DrawLayerItem
	sprite *core.Sprite
}

func NewSprite(parent WorldTransform, im *ebiten.Image, layer Layer) *Sprite {
	w := parent.World()
	e := w.NewEntity()
	spr := &Sprite{}
	spr.WorldItem = newWorldItem(e, w)
	spr.TransformItem = newTransformItem(e, parent)
	spr.DrawLayerItem = newDrawLayerItem(e, w)
	spr.sprite = &core.Sprite{
		ScaleX: 1,
		ScaleY: 1,
		Image:  im,
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawable), spr.sprite); err != nil {
		panic(err)
	}
	return spr
}

func (s *Sprite) SetOffset(x, y float64) {
	s.sprite.OffsetX = x
	s.sprite.OffsetY = y
}

func (s *Sprite) SetOffsetX(x float64) {
	s.sprite.OffsetX = x
}

func (s *Sprite) SetOffsetY(y float64) {
	s.sprite.OffsetY = y
}

func (s *Sprite) SetOrigin(ox, oy float64) {
	s.sprite.OriginX = ox
	s.sprite.OriginY = oy
}

func (s *Sprite) SetImage(img *ebiten.Image) {
	s.sprite.SetImage(img)
}

func (s *Sprite) SetColorTint(c color.Color) {
	s.sprite.SetColorTint(c)
}

// SetColorHue rotates the Hue (in radians)
func (s *Sprite) SetColorHue(theta float64) {
	s.sprite.SetColorHue(theta)
}

func (s *Sprite) ClearColorTransform() {
	s.sprite.ClearColorTransform()
}

func (s *Sprite) SetCompositeMode(mode ebiten.CompositeMode) {
	s.sprite.SetCompositeMode(mode)
}

func init() {
	transparentPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = transparentPixel.Fill(color.Transparent)
}
