package tau

import (
	"testing"

	"github.com/hajimehoshi/ebiten"
)

type c0data struct {
	Points float64
	Scale  float64
}

func TestECS(t *testing.T) {
	s0 := 0.0
	w := NewWorld(&Engine{})
	c0, _ := w.NewComponent(NewComponentInput{
		Name: "test",
	})
	w.NewSystem("", 0, func(ctx Context, screen *ebiten.Image) {
		matches := ctx.System().View().Matches()
		d := matches[0].Components[c0].(*c0data)
		s0 += d.Points * d.Scale * ctx.DT()
	}, c0)
	c0arch := NewArchetype(w, c0)
	c0arch.NewEntity(&c0data{
		Points: 10,
		Scale:  0.25,
	})
	w.Run(nil, 1.0)
	if s0 != 2.5 {
		t.Fail()
	}
}
