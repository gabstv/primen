package utils

import (
	"github.com/hajimehoshi/ebiten"
)

// CloneDrawImageOptions copies the data inside a ebiten.DrawImageOptions
// instance.
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

// CloneAnyDrawImageOptions copies the data of the first non nil
// ebiten.DrawImageOptions instance.
func CloneAnyDrawImageOptions(a ...*ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return CloneDrawImageOptions(v)
		}
	}
	return ebiten.DrawImageOptions{}
}

// AnyDrawImageOptions returns the first non nil *ebiten.DrawImageOptions
func AnyDrawImageOptions(a ...*ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return v
		}
	}
	return &ebiten.DrawImageOptions{}
}
