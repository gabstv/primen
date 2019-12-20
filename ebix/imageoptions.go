package ebix

import (
	"github.com/hajimehoshi/ebiten"
)

func CloneDrawImageOptions(opt *ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	return ebiten.DrawImageOptions{
		GeoM:          opt.GeoM,
		ColorM:        opt.ColorM,
		CompositeMode: opt.CompositeMode,
		Filter:        opt.Filter,
		// ImageParts is deprecated
		// Parts is deprecated
		// SourceRect is deprecated
	}
}

func CloneAnyDrawImageOptions(a ...*ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return CloneDrawImageOptions(v)
		}
	}
	return ebiten.DrawImageOptions{}
}

func AnyDrawImageOptions(a ...*ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return v
		}
	}
	return &ebiten.DrawImageOptions{}
}
