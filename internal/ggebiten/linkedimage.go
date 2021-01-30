package ggebiten

import (
	"image"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

type LinkedImage interface {
	Dispose()
	UpdatePixels()
	Image() *image.RGBA
	Ebimage() *ebiten.Image
}

type HardLinkedImage struct {
	raw   *image.RGBA
	ebimg *ebiten.Image
}

// Dispose is called by GC, use it only if you need to erame memory faster
func (img *HardLinkedImage) Dispose() {
	runtime.SetFinalizer(img, nil)
	img.raw = nil
	img.ebimg.Dispose()
	img.ebimg = nil
}

func (img *HardLinkedImage) discard() {
	runtime.SetFinalizer(img, nil)
	img.raw = nil
	img.ebimg = nil
}

func (img *HardLinkedImage) UpdatePixels() {
	img.ebimg.ReplacePixels(img.raw.Pix)
}

func (img *HardLinkedImage) Image() *image.RGBA {
	return img.raw
}

func (img *HardLinkedImage) Ebimage() *ebiten.Image {
	return img.ebimg
}

type SoftLinkedImage struct {
	raw   *image.RGBA
	ebimg *ebiten.Image
}

// Dispose is called by GC, use it only if you need to erame memory faster
func (img *SoftLinkedImage) Dispose() {
	img.raw = nil
	img.ebimg = nil
}

func (img *SoftLinkedImage) UpdatePixels() {
	img.ebimg.ReplacePixels(img.raw.Pix)
}

func (img *SoftLinkedImage) Image() *image.RGBA {
	return img.raw
}

func (img *SoftLinkedImage) Ebimage() *ebiten.Image {
	return img.ebimg
}

func NewHardLinkedImage(width, height int, filter ebiten.Filter) *HardLinkedImage {
	ebi, _ := ebiten.NewImage(width, height, filter)
	o := &HardLinkedImage{
		raw:   image.NewRGBA(image.Rect(0, 0, width, height)),
		ebimg: ebi,
	}
	runtime.SetFinalizer(o, (*HardLinkedImage).discard)
	return o
}

func NewSoftLinkedImage(width, height int, filter ebiten.Filter) *SoftLinkedImage {
	ebi, _ := ebiten.NewImage(width, height, filter)
	o := &SoftLinkedImage{
		raw:   image.NewRGBA(image.Rect(0, 0, width, height)),
		ebimg: ebi,
	}
	//runtime.SetFinalizer(o, (*SoftLinkedImage).discard)
	return o
}
