package gcs

import (
	"math"
	"testing"

	"github.com/gabstv/ecs"
)

func TestHierarchy(t *testing.T) {
	w := ecs.NewWorld()
	//
	comp := TransformComponent(w)
	TransformSystem(w)
	//
	// the transforms
	tr0 := &Transform{
		X:     100,
		Y:     200,
		Angle: math.Pi / 2,
	}
	//
	tr1 := &Transform{
		X:      10,
		Y:      0,
		Parent: tr0,
	}
	//
	en0 := w.NewEntity()
	w.AddComponentToEntity(en0, comp, tr0)
	en1 := w.NewEntity()
	w.AddComponentToEntity(en1, comp, tr1)
	w.Run(1 / 60)
	x := tr1.globalX // tr1.M.Element(0, 2)
	y := tr1.globalY // tr1.M.Element(1, 2)
	if x != 100 {
		t.Errorf("x != %v (%v) %v", 100, x, tr1.M.String())
	}
	if y != 210 {
		t.Errorf("y != %v (%v)", 210, y)
	}
}
