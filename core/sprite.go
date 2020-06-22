package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// Sprite is the data of a sprite component.
type Sprite struct {
	x        float64 // logical X position
	y        float64 // logical Y position
	angle    float64 // radians
	scaleX   float64 // logical X scale (1 = 100%)
	scaleY   float64 // logical Y scale (1 = 100%)
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
		x:           x,
		y:           y,
		scaleX:      1,
		scaleY:      1,
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

// X gets the local x position. Overrided by the transform
func (s *Sprite) X() float64 {
	return s.x
}

// Y gets the local y position. Overrided by the transform
func (s *Sprite) Y() float64 {
	return s.y
}

func (s *Sprite) SetX(x float64) *Sprite {
	s.x = x
	return s
}

func (s *Sprite) SetY(y float64) *Sprite {
	s.y = y
	return s
}

// Angle gets the local angle (radians).
// It is overrided by the transform component.
func (s *Sprite) Angle() float64 {
	return s.angle
}

func (s *Sprite) SetAngle(r float64) *Sprite {
	s.angle = r
	return s
}

func (s *Sprite) ScaleX() float64 {
	return s.scaleX
}

func (s *Sprite) SetScaleX(sx float64) *Sprite {
	s.scaleX = sx
	return s
}

func (s *Sprite) ScaleY() float64 {
	return s.scaleY
}

func (s *Sprite) SetScaleY(sy float64) *Sprite {
	s.scaleY = sy
	return s
}

func (s *Sprite) Origin() (ox, oy float64) {
	return s.originX, s.originY
}

func (s *Sprite) SetOrigin(ox, oy float64) *Sprite {
	s.originX, s.originY = ox, oy
	return s
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) SetImage(img *ebiten.Image) *Sprite {
	s.image = img
	s.imageWidth, s.imageHeight = getImageSize(img)
	return s
}

func (s *Sprite) Draw(ctx DrawCtx, d *Drawable) {
	if d == nil {
		return
	}
	if s.disabled {
		return
	}
	g := d.G(s.scaleX, s.scaleY, s.angle, s.x, s.y)
	o := &s.opt
	o.GeoM.Reset()
	o.GeoM.Translate(applyOrigin(s.imageWidth, s.originX)+s.offsetX, applyOrigin(s.imageHeight, s.originY)+s.offsetY)
	o.GeoM.Concat(g)

	//TODO: reimplement colormode and composite mode
	ctx.Renderer().DrawImageRaw(s.image, o)

	if DebugDraw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+s.imageWidth, y0
		x2, y2 := x1, y1+s.imageHeight
		x3, y3 := x2-s.imageWidth, y2
		screen := ctx.Renderer().Screen()
		debugLineM(screen, o.GeoM, x0, y0, x1, y1, debugBoundsColor)
		debugLineM(screen, o.GeoM, x1, y1, x2, y2, debugBoundsColor)
		debugLineM(screen, o.GeoM, x2, y2, x3, y3, debugBoundsColor)
		debugLineM(screen, o.GeoM, x3, y3, x0, y0, debugBoundsColor)
		debugLineM(screen, g, -4, 0, 4, 0, debugPivotColor)
		debugLineM(screen, g, 0, -4, 0, 4, debugPivotColor)
	}
}

//go:generate ecsgen -n Sprite -p core -o sprite_component.go --component-tpl --vars "UUID=80C95DEC-DBBF-4529-BD27-739A69055BA0" --vars "BeforeRemove=c.beforeRemove(e)" --vars "OnAdd=c.onAdd(e)"

func (c *SpriteComponent) beforeRemove(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = nil
	}
}

func (c *SpriteComponent) onAdd(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = c.Data(e)
	} else {
		SetDrawableComponentData(c.world, e, Drawable{
			drawer: c.Data(e),
		})
	}
}

//go:generate ecsgen -n DrawableSprite -p core -o sprite_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "OnResize=s.onResize()" --vars "UUID=02F3E9BD-CED3-4160-8943-9A89C0A533FB" --components "Drawable" --components "Sprite"

var matchDrawableSpriteSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetDrawableComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetSpriteComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchDrawableSpriteSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetSpriteComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *DrawableSpriteSystem) onEntityAdded(e ecs.Entity) {
	d := GetDrawableComponentData(s.world, e)
	d.drawer = GetSpriteComponentData(s.world, e)
}

func (s *DrawableSpriteSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *DrawableSpriteSystem) onResize() {
	for _, v := range s.V().Matches() {
		v.Drawable.drawer = v.Sprite
	}
}

// DrawPriority noop
func (s *DrawableSpriteSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop - drawing is controlled by *Drawable
func (s *DrawableSpriteSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *DrawableSpriteSystem) UpdatePriority(ctx UpdateCtx) {}

// Update noop
func (s *DrawableSpriteSystem) Update(ctx UpdateCtx) {}

// func (s *Sprite) SetColorTint(c color.Color) {
// 	s.localColor = ColorTint(c)
// }

// func (s *Sprite) SetColorHue(theta float64) {
// 	s.localColor = ColorM().RotateHue(theta)
// }

// func (s *Sprite) ClearColorTransform() {
// 	s.localColor = nil
// }

// func (s *Sprite) SetCompositeMode(mode ebiten.CompositeMode) {
// 	s.compositeMode = &mode
// }

// // Update does some computation before drawing
// func (s *Sprite) Update(ctx Context) {
// 	if s.lastImage != s.Image {
// 		w, h := s.Image.Size()
// 		s.imageWidth = float64(w)
// 		s.imageHeight = float64(h)
// 		s.lastImage = s.Image
// 	}
// 	if s.localMatrix == nil {
// 		s.localMatrix = GeoM()
// 	}
// 	s.localMatrix.Reset()
// }

// // Draw is called by the Drawable systems
// func (s *Sprite) Draw(renderer DrawManager) {
// 	if s.DrawDisabled {
// 		return
// 	}
// 	g := s.transformMatrix
// 	if g == nil {
// 		g = GeoM().Scale(s.ScaleX, s.ScaleY).Rotate(s.Angle).Translate(s.X, s.Y)
// 	}
// 	s.localMatrix.Translate(applyOrigin(s.imageWidth, s.OriginX), applyOrigin(s.imageHeight, s.OriginY))
// 	s.localMatrix.Translate(s.OffsetX, s.OffsetY)
// 	s.localMatrix.Concat(*g.M())
// 	if s.localColor != nil {
// 		if s.compositeMode != nil {
// 			renderer.DrawImageCComp(s.Image, s.localMatrix, s.localColor, *s.compositeMode)
// 		} else {
// 			renderer.DrawImageC(s.Image, s.localMatrix, s.localColor)
// 		}
// 	} else {
// 		if s.compositeMode != nil {
// 			renderer.DrawImageComp(s.Image, s.localMatrix, *s.compositeMode)
// 		} else {
// 			renderer.DrawImage(s.Image, s.localMatrix)
// 		}
// 	}
// 	if DebugDraw {
// 		x0, y0 := 0.0, 0.0
// 		x1, y1 := x0+s.imageWidth, y0
// 		x2, y2 := x1, y1+s.imageHeight
// 		x3, y3 := x2-s.imageWidth, y2
// 		debugLineM(renderer.Screen(), *s.localMatrix.M(), x0, y0, x1, y1, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *s.localMatrix.M(), x1, y1, x2, y2, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *s.localMatrix.M(), x2, y2, x3, y3, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *s.localMatrix.M(), x3, y3, x0, y0, debugBoundsColor)
// 		debugLineM(renderer.Screen(), *g.M(), -4, 0, 4, 0, debugPivotColor)
// 		debugLineM(renderer.Screen(), *g.M(), 0, -4, 0, 4, debugPivotColor)
// 	}
// }

// func (s *Sprite) Destroy() {
// 	s.Image = nil
// 	s.transformMatrix = nil
// 	s.localMatrix = nil
// }

// func (s *Sprite) IsDisabled() bool {
// 	return s.DrawDisabled
// }
