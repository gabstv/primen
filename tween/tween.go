package tween

import (
	"github.com/gabstv/primen/easing"
)

type Tween struct {
	Function easing.Function
	T        float64
	Setter   func(t01 float64)
	DoneFn   func()
	done     bool
}

func (t *Tween) Update(dt float64) {
	if t.done {
		return
	}
	t.T += dt
	if t.T >= 1 {
		t.T = 1
		t.done = true
	}
	t.Setter(t.Function(t.T))
	if t.done {
		t.DoneFn()
	}
}
