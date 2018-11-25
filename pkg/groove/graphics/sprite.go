package graphics

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove/common"
	"github.com/hajimehoshi/ebiten"
)

var (
	spriteWC *common.WorldComponents
)

func init() {
	spriteWC = &common.WorldComponents{}
}

type SpriteData struct {
	X       float64
	Y       float64
	Options *ebiten.DrawImageOptions
}

func SpriteComponent(w *ecs.World) *ecs.Component {
	return spriteWC.Get(w)
}
