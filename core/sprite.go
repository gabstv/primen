package core

import (
	"image"

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

// Update does some computation before drawing
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

// Draw is called by the Drawable systems
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
	xxg.Translate(s.OffsetX, s.OffsetY)
	xxg.Concat(opt.GeoM)
	centerM := opt.GeoM
	opt.GeoM = *xxg

	if s.lastSubImage != nil {
		screen.DrawImage(s.lastSubImage, opt)
	} else {
		screen.DrawImage(s.Image, opt)
	}
	if DebugDraw {
		x0, y0 := 0.0, 0.0
		x1, y1 := x0+s.imageWidth, y0
		x2, y2 := x1, y1+s.imageHeight
		x3, y3 := x2-s.imageWidth, y2
		debugLineM(screen, opt.GeoM, x0, y0, x1, y1, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x1, y1, x2, y2, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x2, y2, x3, y3, debugBoundsColor)
		debugLineM(screen, opt.GeoM, x3, y3, x0, y0, debugBoundsColor)
		debugLineM(screen, centerM, -4, 0, 4, 0, debugPivotColor)
		debugLineM(screen, centerM, 0, -4, 0, 4, debugPivotColor)
	}
	opt.GeoM = prevGeo
}

// SetTransformMatrix is used by TransformSystem to set a custom transform
func (s *Sprite) SetTransformMatrix(m ebiten.GeoM) {
	s.transformMatrix = m
	s.customMatrix = true
}

func (s *Sprite) ClearTransformMatrix() {
	s.customMatrix = false
}

func (s *Sprite) SetBounds(b image.Rectangle) {
	s.Bounds = b
}

func (s *Sprite) SetOffset(x, y float64) {
	s.OffsetX = x
	s.OffsetY = y
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

// Size returns the real size of the Sprite
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

func (s *Sprite) GetImage() *ebiten.Image {
	return s.Image
}

func (s *Sprite) SetImage(img *ebiten.Image) {
	s.Image = img
}
