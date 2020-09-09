package z

import "github.com/hajimehoshi/ebiten"

func Filter(a string, f ebiten.Filter) ebiten.Filter {
	if a == "" {
		return f
	}
	switch a {
	case "nn", "nearest":
		return ebiten.FilterNearest
	case "linear":
		return ebiten.FilterLinear
	case "default":
		return ebiten.FilterDefault
	}
	return f
}
