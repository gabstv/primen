package main

import (
	"image"
	_ "image/png"
	"math"

	"github.com/gabstv/troupe"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)

	engine := troupe.NewEngine(&troupe.NewEngineInput{
		Title:  "Basic Sprites",
		Width:  320,
		Height: 240,
		Scale:  2,
	})

	//go func(){
	dw := engine.Default()
	sc := troupe.SpriteComponent(dw)
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, sc, &troupe.Sprite{
		Image:  ebimg,
		X:      64,
		Y:      64,
		Angle:  math.Pi / 2,
		ScaleX: 1,
		ScaleY: 1,
		Bounds: image.Rect(0, 0, 16, 16),
	})
	e2 := dw.NewEntity()
	dw.AddComponentToEntity(e2, sc, &troupe.Sprite{
		Image:   ebimg,
		X:       128,
		Y:       64,
		Angle:   0,
		ScaleX:  2,
		ScaleY:  2,
		OriginX: -.5,
		OriginY: -.5,
		Bounds:  image.Rect(16, 16, 32, 32),
	})
	dw.NewSystem("", 0, func(ctx troupe.Context, screen *ebiten.Image) {
		matches := ctx.System().View().Matches()
		dt := ctx.DT()
		for _, m := range matches {
			sprite := m.Components[sc].(*troupe.Sprite)
			sprite.Angle = sprite.Angle + (math.Pi * dt * 0.0125 * 4)
		}
	}, sc)
	//}()
	engine.Run()
}
