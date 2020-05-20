package core

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

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
