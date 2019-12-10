package ebix

import (
	"github.com/hajimehoshi/ebiten"
)

func CloneDrawImageOptions(opt *ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	return ebiten.DrawImageOptions{
		GeoM: opt.GeoM,
		ColorM: opt.ColorM,
		CompositeMode: opt.CompositeMode,
		Filter: opt.Filter,
		// ImageParts is deprecated
		// Parts is deprecated
		// SourceRect is deprecated
	}
}