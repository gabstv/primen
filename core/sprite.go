package core

import (
	"github.com/hajimehoshi/ebiten"
)

// Sprite is the data of a sprite component.
type Sprite struct {
	X       float64 // logical X position
	Y       float64 // logical Y position
	Angle   float64 // radians
	ScaleX  float64 // logical X scale (1 = 100%)
	ScaleY  float64 // logical Y scale (1 = 100%)
	OriginX float64 // X origin (0 = left; 0.5 = center; 1 = right)
	OriginY float64 // Y origin (0 = top; 0.5 = middle; 1 = bottom)
	OffsetX float64 // offset origin X (in pixels)
	OffsetY float64 // offset origin Y (in pixels)
	Image   *ebiten.Image

	DrawDisabled bool // if true, the SpriteSystem will not draw this

	lastImage *ebiten.Image // lastImage exists to keep track of the public Image field, if it
	// changes, the imageWidth and ImageHeight needs to be recalculated.
	imageWidth  float64 // last calculated image width
	imageHeight float64 // last calculated image height
	//
	transformMatrix GeoMatrix
	localMatrix     GeoMatrix
}

// Update does some computation before drawing
func (s *Sprite) Update(ctx Context) {
	if s.lastImage != s.Image {
		w, h := s.Image.Size()
		s.imageWidth = float64(w)
		s.imageHeight = float64(h)
		s.lastImage = s.Image
	}
	if s.localMatrix == nil {
		s.localMatrix = GeoM()
	}
	s.localMatrix.Reset()
}

// Draw is called by the Drawable systems
func (s *Sprite) Draw(renderer DrawManager) {
	if s.DrawDisabled {
		return
	}
	g := s.transformMatrix
	if g == nil {
		g = GeoM().Scale(s.ScaleX, s.ScaleY).Rotate(s.Angle).Translate(s.X, s.Y)
	}
	s.localMatrix.Translate(applyOrigin(s.imageWidth, s.OriginX), applyOrigin(s.imageHeight, s.OriginY))
	s.localMatrix.Translate(s.OffsetX, s.OffsetY)
	s.localMatrix.Concat(*g.M())
	renderer.DrawImageG(s.Image, s.localMatrix)
	if DebugDraw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+s.imageWidth, y0
		x2, y2 := x1, y1+s.imageHeight
		x3, y3 := x2-s.imageWidth, y2
		debugLineM(renderer.Screen(), *s.localMatrix.M(), x0, y0, x1, y1, debugBoundsColor)
		debugLineM(renderer.Screen(), *s.localMatrix.M(), x1, y1, x2, y2, debugBoundsColor)
		debugLineM(renderer.Screen(), *s.localMatrix.M(), x2, y2, x3, y3, debugBoundsColor)
		debugLineM(renderer.Screen(), *s.localMatrix.M(), x3, y3, x0, y0, debugBoundsColor)
		debugLineM(renderer.Screen(), *g.M(), -4, 0, 4, 0, debugPivotColor)
		debugLineM(renderer.Screen(), *g.M(), 0, -4, 0, 4, debugPivotColor)
	}
}

// SetTransformMatrix is used by TransformSystem to set a custom transform
func (s *Sprite) SetTransformMatrix(m GeoMatrix) {
	s.transformMatrix = m
}

func (s *Sprite) ClearTransformMatrix() {
	s.transformMatrix = nil
}

func (s *Sprite) SetOffset(x, y float64) {
	s.OffsetX = x
	s.OffsetY = y
}

func (s *Sprite) Destroy() {
	s.Image = nil
	s.transformMatrix = nil
	s.localMatrix = nil
}

func (s *Sprite) IsDisabled() bool {
	return s.DrawDisabled
}

// Size returns the real size of the Sprite
func (s *Sprite) Size() (w, h float64) {
	return s.imageWidth, s.imageHeight
}

// GetPrecomputedImage returns the last precomputed image
//
// TODO: remove
func (s *Sprite) GetPrecomputedImage() *ebiten.Image {
	return s.lastImage
}

// GetPrecomputedImageDim returns the last precomputed image dimmensions
func (s *Sprite) GetPrecomputedImageDim() (width, height float64) {
	return s.imageWidth, s.imageHeight
}

func (s *Sprite) GetImage() *ebiten.Image {
	return s.Image
}

func (s *Sprite) SetImage(img *ebiten.Image) {
	s.Image = img
}

var _ Drawable = &Sprite{}
