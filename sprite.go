package primen

import (
	"image"
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

type Animation = core.Animation

type AnimatedSprite struct {
	*Sprite
	coreAnim *core.SpriteAnimation
}

func NewAnimatedSprite(parent WorldTransform, layer Layer, anim Animation) *AnimatedSprite {
	as := &AnimatedSprite{
		Sprite: NewSprite(parent, transparentPixel, layer),
	}
	sa := &core.SpriteAnimation{
		Enabled: true,
		Anim:    anim,
	}
	if err := as.World().AddComponentToEntity(as.Entity(), as.World().Component(core.CNSpriteAnimation), sa); err != nil {
		panic(err)
	}
	as.coreAnim = sa
	return as
}

func (as *AnimatedSprite) PlayClipIndex(i int) {
	as.coreAnim.PlayClipIndex(i)
}

func (as *AnimatedSprite) PlayClip(name string) {
	as.coreAnim.PlayClip(name)
}

type TileSet struct {
	*WorldItem
	*TransformItem
	*DrawLayerItem
	tileset *core.TileSet
}

func (t *TileSet) SetDB(db []*ebiten.Image) {
	t.tileset.DB = db
}

func (t *TileSet) SetTilesXY(m2d [][]int) {
	if len(m2d) < 1 {
		return
	}
	width := len(m2d)
	height := len(m2d[0])
	mm := make([]int, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			mm[y*width+x] = m2d[x][y]
		}
	}
	t.tileset.Cells = mm
	t.tileset.CSize = image.Point{
		X: width,
		Y: height,
	}
}

func (t *TileSet) SetTilesYX(m2d [][]int) {
	if len(m2d) < 1 {
		return
	}
	height := len(m2d)
	width := len(m2d[0])
	mm := make([]int, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mm[y*width+x] = m2d[y][x]
		}
	}
	t.tileset.Cells = mm
	t.tileset.CSize = image.Point{
		X: width,
		Y: height,
	}
}

func (t *TileSet) SetCellSize(w, h float64) {
	t.tileset.CellWidth = w
	t.tileset.CellHeight = h
}

func (t *TileSet) SetOrigin(ox, oy float64) {
	t.tileset.OriginX = ox
	t.tileset.OriginY = oy
}

func NewTileSet(parent WorldTransform, layer Layer) *TileSet {
	w := parent.World()
	e := w.NewEntity()
	spr := &TileSet{}
	spr.WorldItem = newWorldItem(e, w)
	spr.TransformItem = newTransformItem(e, parent)
	spr.DrawLayerItem = newDrawLayerItem(e, w)
	spr.tileset = &core.TileSet{
		ScaleX: 1,
		ScaleY: 1,
	}
	if err := w.AddComponentToEntity(spr.entity, w.Component(core.CNDrawable), spr.tileset); err != nil {
		panic(err)
	}
	return spr
}

func init() {
	transparentPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = transparentPixel.Fill(color.Transparent)
}
