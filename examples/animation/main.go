package main

// https://www.kenney.nl/assets/platformer-characters-1

import (
	"math"
	_ "image/png"
	"image"
	"fmt"
	
	"github.com/gabstv/ecs"
	"github.com/gabstv/troupe/pkg/troupe"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("zombie_kenney.png", ebiten.FilterDefault)
	ppimg, _, _ := ebitenutil.NewImageFromFile("ping_pong.png", ebiten.FilterDefault)

	engine := troupe.NewEngine(&troupe.NewEngineInput{
		Title: "Basic Animation",
		Width: 640,
		Height: 480,
		Scale: 1,
	})
	
	dw := engine.Default()
	sc := troupe.SpriteComponent(dw)
	ac := troupe.SpriteAnimationComponent(dw)
	createCharacter(dw, sc, ac, ebimg)
	createPingPonger(dw, sc, ac, ppimg)

	s0 := dw.NewSystem(0, func(dt float64, v *ecs.View, s *ecs.System){
		fps := ebiten.CurrentFPS()
		img := engine.Get(troupe.EbitenScreen).(*ebiten.Image)
		ebitenutil.DebugPrintAt(img, fmt.Sprintf("%.2f fps", fps), 0, 0)
	})
	s0.AddTag(troupe.WorldTagDraw)
	
	engine.Run()
}

func createCharacter(dw *ecs.World, spriteComp *ecs.Component, animComp *ecs.Component, ebimg *ebiten.Image){
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, spriteComp, &troupe.Sprite{
		Image: ebimg,
		X: 300,
		Y: 200,
		Angle: 0,
		ScaleX: 1,
		ScaleY: 1,
		Bounds: image.Rect(0,0,80,110),
	})
	dw.AddComponentToEntity(e, animComp, &troupe.SpriteAnimation{
		Enabled: true,
		Playing :true,
		Clips: []troupe.SpriteAnimationClip{
			troupe.SpriteAnimationClip{
				Name: "default",
				Frames: []image.Rectangle{
					image.Rect(0,0,80,110), // 0
					image.Rect(80,0,80*2,110), // 1
					image.Rect(80*2,0,80*3,110), // 2
					image.Rect(80*3,0,80*4,110), // 3
					image.Rect(0,0,80,110), // 0
				},
				ClipMode: troupe.AnimLoop,
			},
		},
		Fps: 24,
	})
}

func createPingPonger(dw *ecs.World, spriteComp *ecs.Component, animComp *ecs.Component, ebimg *ebiten.Image){
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, spriteComp, &troupe.Sprite{
		Image: ebimg,
		X: 370,
		Y: 180,
		Angle: math.Pi / 4,
		ScaleX: 1,
		ScaleY: 1,
		Bounds: image.Rect(0,0,8,32),
	})
	dw.AddComponentToEntity(e, animComp, &troupe.SpriteAnimation{
		Enabled: true,
		Playing :true,
		Clips: []troupe.SpriteAnimationClip{
			troupe.SpriteAnimationClip{
				Name: "default",
				Frames: []image.Rectangle{
					image.Rect(8*0,0,8*1,32),
					image.Rect(8*1,0,8*2,32),
					image.Rect(8*2,0,8*3,32),
					image.Rect(8*3,0,8*4,32),
				},
				ClipMode: troupe.AnimPingPong,
			},
		},
		Fps: 24,
	})
}