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
	xy := tr1.M.Project(ZV) // tr1.M.Element(0, 2)
	if xy.X != 100 {
		t.Errorf("x != %v (%v)", 100, xy.X)
	}
	if xy.Y != 210 {
		t.Errorf("y != %v (%v)", 210, xy.Y)
	}
}

func TestMatrices(t *testing.T) {
	mA := IM
	mA = mA.Rotated(ZV, math.Pi)
	mB := IM
	mB = mB.Moved(V(0.5, 0))
	mC := mA.Chained(mB)
	if mC[0] != -1 {
		t.FailNow()
	}
	v2 := mC.Project(V(-1, 1))
	if v2.X != 1.5 {
		t.Fatalf("v2.X != 1.5 (%v)", v2.X)
	}
	if v2.Y != -1 {
		t.Fatalf("v2.Y != -1 (%v)", v2.Y)
	}
}
