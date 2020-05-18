package main

import (
	"image"
	_ "image/png"
	"math"

	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var anglecs = &tau.BasicCS{
	SysName: "anglecs",
	SysExec: func(ctx tau.Context) {
		sc := ctx.World().Component(tau.CNDrawable)
		matches := ctx.System().View().Matches()
		dt := ctx.DT()
		for _, m := range matches {
			sprite := m.Components[sc].(*tau.Sprite)
			sprite.Angle = sprite.Angle + (math.Pi * dt * 0.0125 * 4)
		}
	},
	GetComponents: func(w *ecs.World) []*ecs.Component {
		return []*ecs.Component{
			w.Component(tau.CNDrawable),
		}
	},
}

func main() {
	ebimg, _, err := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	engine := tau.NewEngine(&tau.NewEngineInput{
		Title:  "Basic Sprites",
		Width:  320,
		Height: 240,
		Scale:  2,
	})

	dw := engine.Default()
	sc := dw.Component(tau.CNDrawable)
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, sc, &tau.Sprite{
		Image:   ebimg,
		X:       64,
		Y:       64,
		Angle:   math.Pi / 2,
		ScaleX:  1,
		ScaleY:  1,
		OriginX: 0.5,
		OriginY: 0.5,
		Bounds:  image.Rect(0, 0, 16, 16),
	})
	e2 := dw.NewEntity()
	dw.AddComponentToEntity(e2, sc, &tau.Sprite{
		Image:   ebimg,
		X:       128,
		Y:       64,
		Angle:   0,
		ScaleX:  2,
		ScaleY:  2,
		OriginX: 0,
		OriginY: 0,
		Bounds:  image.Rect(16, 16, 32, 32),
	})

	tau.SetupSystem(dw, anglecs)

	engine.Run()
}
