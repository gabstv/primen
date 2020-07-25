package pb

import (
	"github.com/gabstv/ebiten"
)

func ToEbitenFilter(t ImageFilter) ebiten.Filter {
	switch t {
	case ImageFilter_LINEAR:
		return ebiten.FilterLinear
	case ImageFilter_NEAREST:
		return ebiten.FilterNearest
	}
	return 0
}
