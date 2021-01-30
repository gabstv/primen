package z

import "github.com/hajimehoshi/ebiten/v2"

func Filter(a string, f ebiten.Filter) ebiten.Filter {
	if a == "" {
		return f
	}
	switch a {
	case "nn", "nearest":
		return ebiten.FilterNearest
	case "linear":
		return ebiten.FilterLinear
	}
	return f
}
