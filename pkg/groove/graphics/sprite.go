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

// Sprite is the data of a sprite component.
type Sprite struct {
	X       float64
	Y       float64
	Options *ebiten.DrawImageOptions
}

// SpriteComponent will get the registered sprite component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func SpriteComponent(w *ecs.World) *ecs.Component {
	c := spriteWC.Get(w)
	if c == nil {
		var err error
		c, err = w.NewComponent(ecs.NewComponentInput{
			Name: "groove.graphics.Sprite",
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*Sprite)
				return ok
			},
			DestructorFn: func(w2 *ecs.World, entity ecs.Entity, data interface{}) {
				sd := data.(*Sprite)
				sd.Options = nil
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return c
}
