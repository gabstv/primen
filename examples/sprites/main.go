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
		sc := gcs.SpriteComponent(dw)
		e := dw.NewEntity()
		dw.AddComponentToEntity(e, sc, &gcs.Sprite{
			Image: ebimg,
			X: 64,
			Y: 64,
			Angle: math.Pi/2,
			ScaleX: 1,
			ScaleY: 1,
			Bounds: image.Rect(0,0,16,16),
		})
		e2 := dw.NewEntity()
		dw.AddComponentToEntity(e2, sc, &gcs.Sprite{
			Image: ebimg,
			X: 128,
			Y: 64,
			Angle: 0,
			ScaleX: 2,
			ScaleY: 2,
			Bounds: image.Rect(16,16,32,32),
		})
		dw.NewSystem(0, func(dt float64, v *ecs.View, s *ecs.System){
			matches := v.Matches()
			for _, m := range matches {
				sprite := m.Components[sc].(*gcs.Sprite)
				sprite.Angle = sprite.Angle + (math.Pi * dt * 0.0125)
			}
		}, sc)
	//}()
	engine.Run()
}
