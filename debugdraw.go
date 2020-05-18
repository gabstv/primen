package tau

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	debugBoundsColor = color.RGBA{
		R: 255,
		B: 255,
		A: 255,
	}
)

var debugPixel *ebiten.Image

func colorScale(clr color.Color) (rf, gf, bf, af float64) {
	r, g, b, a := clr.RGBA()
	if a == 0 {
		return 0, 0, 0, 0
	}

	rf = float64(r) / float64(a)
	gf = float64(g) / float64(a)
	bf = float64(b) / float64(a)
	af = float64(a) / 0xffff
	return
}

func debugLineM(dst *ebiten.Image, srcopt *ebiten.DrawImageOptions, x1, y1, x2, y2 float64, clr color.Color) {
	xx1, yx1 := srcopt.GeoM.Apply(x1, y1)
	xx2, yx2 := srcopt.GeoM.Apply(x2, y2)
	ebitenutil.DrawLine(dst, xx1, yx1, xx2, yx2, clr)
}

func debugLineM2(dst *ebiten.Image, srcopt *ebiten.DrawImageOptions, x1, y1, x2, y2 float64, clr color.Color) {
	srcColor := srcopt.ColorM
	srcGeo := srcopt.GeoM

	//srcopt.GeoM.Apply(x1, y1)
	//
	/*
		ew, eh := emptyImage.Size()
		length := math.Hypot(x2-x1, y2-y1)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(length/float64(ew), 1/float64(eh))
		op.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
		op.GeoM.Translate(x1, y1)
		op.ColorM.Scale(colorScale(clr))
		// Filter must be 'nearest' filter (default).
		// Linear filtering would make edges blurred.
		_ = dst.DrawImage(emptyImage, op)
	*/
	//
	//
	ew, eh := debugPixel.Size()
	length := math.Hypot(x2-x1, y2-y1)
	mm := &ebiten.GeoM{}
	mm.Scale(length/float64(ew), 1/float64(eh))
	mm.Rotate(math.Atan2(y2-y1, x2-x1))
	mm.Translate(x1, y1)
	mm.Concat(srcGeo)
	srcopt.GeoM = *mm
	//srcopt.ColorM.Scale(colorScale(clr))
	dst.DrawImage(debugPixel, srcopt)
	srcopt.GeoM = srcGeo
	srcopt.ColorM = srcColor
}

func init() {
	debugPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = debugPixel.Fill(color.White)
}
