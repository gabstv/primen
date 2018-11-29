package main

// https://www.kenney.nl/assets/platformer-characters-1

import (
	_ "image/png"
	"image"
	"fmt"
	
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/gcs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("zombie_kenney.png", ebiten.FilterDefault)

	engine := groove.NewEngine(&groove.NewEngineInput{
		Title: "Basic Animation",
		Width: 640,
		Height: 480,
		Scale: 1,
	})
	
	dw := engine.Default()
	sc := gcs.SpriteComponent(dw)
	ac := gcs.SpriteAnimationComponent(dw)
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, sc, &gcs.Sprite{
		Image: ebimg,
		X: 300,
		Y: 200,
		Angle: 0,
		ScaleX: 1,
		ScaleY: 1,
		Bounds: image.Rect(0,0,80,110),
	})
	dw.AddComponentToEntity(e, ac, &gcs.SpriteAnimation{
		Enabled: true,
		Play :true,
		Clips: []gcs.SpriteAnimationClip{
			gcs.SpriteAnimationClip{
				Name: "default",
				Frames: []image.Rectangle{
					image.Rect(0,0,80,110), // 0
					image.Rect(80,0,80*2,110), // 1
					image.Rect(80*2,0,80*3,110), // 2
					image.Rect(80*3,0,80*4,110), // 3
					image.Rect(0,0,80,110), // 0
				},
				ClipMode: gcs.AnimLoop,
			},
		},
		Fps: 24,
	})

	s0 := dw.NewSystem(0, func(dt float64, v *ecs.View, s *ecs.System){
		fps := ebiten.CurrentFPS()
		img := engine.Get(groove.EbitenScreen).(*ebiten.Image)
		ebitenutil.DebugPrintAt(img, fmt.Sprintf("%.2f fps", fps), 0, 0)
	})
	s0.AddTag(groove.WorldTagDraw)
	
	engine.Run()
}
