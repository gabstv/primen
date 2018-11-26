package main

import (
	"github.com/gabstv/ecs"
	"math"
	_ "image/png"
	"image"

	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/graphics"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)

	engine := groove.NewEngine(&groove.NewEngineInput{
		Title: "Basic Sprites",
		Width: 320,
		Height: 240,
		Scale: 2,
	})
	
	//go func(){
		dw := engine.Default()
		sc := graphics.SpriteComponent(dw)
		e := dw.NewEntity()
		dw.AddComponentToEntity(e, sc, &graphics.Sprite{
			Image: ebimg,
			X: 64,
			Y: 64,
			Angle: math.Pi/2,
			Bounds: image.Rect(0,0,16,16),
		})
		e2 := dw.NewEntity()
		dw.AddComponentToEntity(e2, sc, &graphics.Sprite{
			Image: ebimg,
			X: 128,
			Y: 64,
			Angle: math.Pi/4,
			Bounds: image.Rect(16,16,32,32),
		})
		dw.NewSystem(0, func(dt float64, v *ecs.View){
			matches := v.Matches()
			for _, m := range matches {
				sprite := m.Components[sc].(*graphics.Sprite)
				sprite.Angle = sprite.Angle + (math.Pi * dt)
			}
		}, sc)
	//}()
	engine.Run()
}
