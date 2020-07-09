package ggui

import "github.com/gabstv/ecs/v2"

// SizeableUI is the component that controls the width and height of a UI
type SizeableUI interface {
	SetSize(w, h float64)
}

type GetSizeableUIFn func(w ecs.BaseWorld, e ecs.Entity) SizeableUI

func RegisterSizeableUIComponent(w ecs.BaseWorld, f ecs.Flag, fn GetSizeableUIFn) {
	gflags := w.FlagGroup("PrimenSizeableUI")
	gflags = gflags.Or(f)
	w.SetFlagGroup("PrimenSizeableUI", gflags)
	vi := w.LGet("PrimenSizeableUI")
	var vs map[uint8]GetSizeableUIFn
	if vi != nil {
		vs = vi.(map[uint8]GetSizeableUIFn)
	} else {
		vs = make(map[uint8]GetSizeableUIFn)
	}
	vs[f.Lowest()] = fn
	w.LSet("PrimenSizeableUI", vs)
}

func GetSizeableUI(w ecs.BaseWorld, e ecs.Entity) SizeableUI {
	eflag := w.CFlag(e)
	dflag := eflag.And(w.FlagGroup("PrimenSizeableUI"))
	vi := w.LGet("PrimenSizeableUI").(map[uint8]GetSizeableUIFn)
	getter := vi[dflag.Lowest()]
	return getter(w, e)
}
