// +build example

package main

import (
	"image"
	_ "image/png"
	"math"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var anglecs = &core.BasicCS{
	SysName: "anglecs",
	SysExec: func(ctx core.Context) {
		sc := ctx.World().Component(core.CNDrawable)
		matches := ctx.System().View().Matches()
		dt := ctx.DT()
		for _, m := range matches {
			sprite := m.Components[sc].(*core.Sprite)
			sprite.Angle = sprite.Angle + (math.Pi * dt * 0.0125 * 4)
		}
	},
	GetComponents: func(w *ecs.World) []*ecs.Component {
		return []*ecs.Component{
			w.Component(core.CNDrawable),
		}
	},
}

func main() {
	ebimg, _, err := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	engine := primen.NewEngine(&primen.NewEngineInput{
		Title:  "Basic Sprites",
		Width:  640,
		Height: 480,
		Scale:  2,
	})

	dw := engine.Default()
	sc := dw.Component(core.CNDrawable)
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, sc, &core.Sprite{
		Image:   ebimg.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		X:       64,
		Y:       64,
		Angle:   math.Pi / 2,
		ScaleX:  1,
		ScaleY:  1,
		OriginX: 0.5,
		OriginY: 0.5,
	})
	e2 := dw.NewEntity()
	dw.AddComponentToEntity(e2, sc, &core.Sprite{
		Image:   ebimg.SubImage(image.Rect(16, 16, 32, 32)).(*ebiten.Image),
		X:       128,
		Y:       64,
		Angle:   0,
		ScaleX:  2,
		ScaleY:  2,
		OriginX: 0,
		OriginY: 0,
	})

	core.SetupSystem(dw, anglecs)

	engine.Run()
}
