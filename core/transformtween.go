package core

import (
	"sort"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/easing"
)

type TrTweenType int

const (
	TrTweenX        TrTweenType = 1
	TrTweenY        TrTweenType = 2
	TrTweenRotation TrTweenType = 3
	TrTweenScaleXY  TrTweenType = 4
	TrTweenScaleX   TrTweenType = 5
	TrTweenScaleY   TrTweenType = 6
)

type TrTween struct {
	Easing   easing.Function
	From     float64
	To       float64
	Duration float64
	Type     TrTweenType
}

type TrTweenTuple struct {
	Name  string
	Tween TrTween
}

type TrTweening struct {
	tweens     []TrTweenTuple
	disabled   bool
	active     TrTween
	activename string
	t          float64
	playing    bool
	callback   func(name string)
}

func NewTrTweening() TrTweening {
	return TrTweening{
		tweens:  make([]TrTweenTuple, 0, 1),
		playing: false,
	}
}

func (t *TrTweening) Play(name string) bool {
	i := sort.Search(len(t.tweens), func(i int) bool {
		return t.tweens[i].Name >= name
	})
	if i >= 0 && i < len(t.tweens) && t.tweens[i].Name == name {
		t.t = 0
		t.active = t.tweens[i].Tween
		t.activename = name
		t.playing = true
		return true
	}
	return false
}

func (t *TrTweening) SetTween(name string, ttype TrTweenType, from, to, duration float64, easingfn easing.Function) *TrTweening {
	if duration <= 0 {
		duration = 1
	}
	if t.tweens == nil {
		t.tweens = make([]TrTweenTuple, 0)
	}
	i := sort.Search(len(t.tweens), func(i int) bool {
		return t.tweens[i].Name >= name
	})
	if i >= 0 && i < len(t.tweens) && t.tweens[i].Name == name {
		t.tweens[i] = TrTweenTuple{
			Name: name,
			Tween: TrTween{
				Easing:   easingfn,
				From:     from,
				To:       to,
				Duration: duration,
				Type:     ttype,
			},
		}
		return t
	}
	t.tweens = append(t.tweens, TrTweenTuple{
		Name: name,
		Tween: TrTween{
			Easing:   easingfn,
			From:     from,
			To:       to,
			Duration: duration,
			Type:     ttype,
		},
	})
	sort.Slice(t.tweens, func(i, j int) bool {
		return t.tweens[i].Name < t.tweens[j].Name
	})
	return t
}

func (t *TrTweening) RemoveTween(name string) bool {
	if t.tweens == nil {
		t.tweens = make([]TrTweenTuple, 0)
	}
	i := sort.Search(len(t.tweens), func(i int) bool {
		return t.tweens[i].Name >= name
	})
	if i >= 0 && i < len(t.tweens) && t.tweens[i].Name == name {
		t.tweens = t.tweens[:i+copy(t.tweens[i:], t.tweens[i+1:])]
		return true
	}
	return false
}

func (t *TrTweening) SetDoneCallback(fn func(name string)) *TrTweening {
	t.callback = fn
	return t
}

//go:generate ecsgen -n TrTweening -p core -o trtweening_component.go --component-tpl --vars "UUID=7D0BCDA8-ABE8-41EB-BF23-7DDCB4152AFD"

//go:generate ecsgen -n TrTweening -p core -o trtweening_system.go --system-tpl --vars "Priority=90" --vars "UUID=820C75AB-CAD6-47AE-A84C-1EC7BAECE328" --components "Transform" --components "TrTweening"

var matchTrTweeningSystem = func(eflag ecs.Flag, world ecs.BaseWorld) bool {
	return eflag.Contains(GetTransformComponent(world).Flag().Or(GetTrTweeningComponent(world).Flag()))
}

var resizematchTrTweeningSystem = func(eflag ecs.Flag, world ecs.BaseWorld) bool {
	return eflag.ContainsAny(GetTransformComponent(world).Flag().Or(GetTrTweeningComponent(world).Flag()))
}

// DrawPriority noop
func (s *TrTweeningSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop
func (s *TrTweeningSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *TrTweeningSystem) UpdatePriority(ctx UpdateCtx) {}

// Update resolves active tweens
func (s *TrTweeningSystem) Update(ctx UpdateCtx) {
	dt := ctx.DT()
	for _, v := range s.V().Matches() {
		if v.TrTweening.disabled {
			continue
		}
		if !v.TrTweening.playing {
			continue
		}
		done := false
		v.TrTweening.t += dt / v.TrTweening.active.Duration
		if v.TrTweening.t >= 1.0 {
			v.TrTweening.t = 1.0
			done = true
		}
		vt := v.TrTweening.active.Easing(v.TrTweening.t)
		vf := Lerpf(v.TrTweening.active.From, v.TrTweening.active.To, vt)
		switch v.TrTweening.active.Type {
		case TrTweenX:
			v.Transform.SetX(vf)
		case TrTweenY:
			v.Transform.SetY(vf)
		case TrTweenRotation:
			v.Transform.SetAngle(vf)
		case TrTweenScaleXY:
			v.Transform.SetScale(vf, vf)
		case TrTweenScaleX:
			v.Transform.SetScaleX(vf)
		case TrTweenScaleY:
			v.Transform.SetScaleY(vf)
		}
		if done {
			v.TrTweening.playing = false
			if v.TrTweening.callback != nil {
				v.TrTweening.callback(v.TrTweening.activename)
			}
		}
	}
}
