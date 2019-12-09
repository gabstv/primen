package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten"
	"github.com/stretchr/testify/assert"
)

type compAData struct {
	X int
	Y int
}
type compBData struct {
	Name string
}

func TestNewArchetype(t *testing.T) {
	w := NewWorld()

	c1, err := w.NewComponent(NewComponentInput{
		Name: "COMP_A",
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*compAData)
			return ok
		},
	})
	assert.NoError(t, err)
	c2, err := w.NewComponent(NewComponentInput{
		Name: "COMP_B",
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*compBData)
			return ok
		},
	})
	assert.NoError(t, err)
	arche1 := NewArchetype(w, c1, c2)

	w.NewSystem(1, func(ctx Context, screen *ebiten.Image) {
		m := ctx.System().View().Matches()
		for _, v := range m {
			da := v.Components[c1].(*compAData)
			db := v.Components[c2].(*compBData)
			if db.Name == "Troupe" {
				da.X++
			} else {
				da.Y++
			}
		}
	}, c1, c2)

	// most optimal way to instantiate an archetype is to follow the order
	// of the components specified by NewArchetype
	e1 := arche1.NewEntity(&compAData{
		X: 10,
		Y: 20,
	}, &compBData{
		Name: "Troupe",
	})
	e2 := arche1.NewEntity(&compAData{
		X: 10,
		Y: 20,
	}, &compBData{
		Name: "Trends",
	})
	w.Run(nil, 1)

	ed1 := c1.Data(e1).(*compAData)
	ed2 := c1.Data(e2).(*compAData)

	assert.Equal(t, 11, ed1.X)
	assert.Equal(t, 20, ed1.Y)
	assert.Equal(t, 10, ed2.X)
	assert.Equal(t, 21, ed2.Y)
}
