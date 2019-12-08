package troupe

import (
	"math"
	"testing"
)

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
	if math.Abs(v2.X-1.5) > 0.00001 {
		t.Fatalf("v2.X != 1.5 (%v)", v2.X)
	}
	if math.Abs(v2.Y-(-1)) > 0.00001 {
		t.Fatalf("v2.Y != -1 (%v)", v2.Y)
	}
}
