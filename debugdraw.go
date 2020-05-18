package tau

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	debugBoundsColor = color.RGBA{
		R: 255,
		B: 255,
		A: 230,
	}
	debugPivotColor = color.RGBA{
		R: 255,
		A: 230,
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

func debugLineM(dst *ebiten.Image, m ebiten.GeoM, x1, y1, x2, y2 float64, clr color.Color) {
	xx1, yx1 := m.Apply(x1, y1)
	xx2, yx2 := m.Apply(x2, y2)
	ebitenutil.DrawLine(dst, xx1, yx1, xx2, yx2, clr)
}

func init() {
	debugPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = debugPixel.Fill(color.White)
}
