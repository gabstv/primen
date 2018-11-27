package main

import (
	"math"
	_ "image/png"
	"image"
	
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/gcs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type spinner struct{
	Speed float64
}

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)

	engine := groove.NewEngine(&groove.NewEngineInput{
		Title: "Basic Transform With Sprites",
		Width: 320,
		Height: 240,
		Scale: 2,
	})
	
	dw := engine.Default()
	sc := gcs.SpriteComponent(dw)
	tc := gcs.TransformComponent(dw)
	spinnercomp, _ := dw.NewComponent(ecs.NewComponentInput{
		Name: "spinner",
	})
	ss := dw.NewSystem(1, spinnersys, tc, spinnercomp)
	ss.Set("spinnercomp", spinnercomp)
	ss.Set("tc", tc)
	e := dw.NewEntity()
	t99 := &gcs.Transform{
		X: 320/2,
		Y: 240/2,
	}
	dw.AddComponentToEntity(e, tc, t99)
	dw.AddComponentToEntity(e, spinnercomp, &spinner{
		Speed: 1,
	})
	// add children
	for i := 0; i < 10; i++ {
		e2 := dw.NewEntity()
		mm := gcs.IM.Moved(gcs.V(10,0)).Rotated(gcs.ZV, (math.Pi*2)/float64(10-i)).Project(gcs.ZV)
		println(mm.String())
		dw.AddComponentToEntity(e2, tc, &gcs.Transform{
			X: mm.X,
			Y: mm.Y,
			Parent: t99,
		})
		dw.AddComponentToEntity(e2, sc, &gcs.Sprite{
			Bounds: image.Rect(0,0,16,16),
			Image: ebimg,
			ScaleX: 1,
			ScaleY:1,
		})
	}
	engine.Run()
}

func spinnersys(dt float64, view *ecs.View, sys *ecs.System) {
	sc := sys.Get("spinnercomp").(*ecs.Component)
	tc := sys.Get("tc").(*ecs.Component)
	//
	for _, v := range view.Matches() {
		spin := v.Components[sc].(*spinner)
		tr := v.Components[tc].(*gcs.Transform)
		tr.Angle += spin.Speed*dt
	}
}