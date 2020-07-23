// +build example

package main

// https://www.kenney.nl/assets/platformer-characters-1

import (
	"image"
	_ "image/png"
	"math"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components/graphics"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Title:   "Basic Animation",
		Width:   640,
		Height:  480,
		Scale:   1,
		OnReady: einit,
	})

	//core.DebugDraw = true
	engine.SetDebugTPS(true)

	engine.Run()
}

func einit(engine primen.Engine) {
	ebimg, _, _ := ebitenutil.NewImageFromFile("zombie_kenney.png", ebiten.FilterDefault)
	ppimg, _, _ := ebitenutil.NewImageFromFile("ping_pong.png", ebiten.FilterDefault)
	dw := engine.NewWorldWithDefaults(0)
	createCharacter(dw, ebimg)
	createPingPonger(dw, ppimg, 370, 180, 0)
	createPingPonger(dw, ppimg, 380, 180, 1)
	createPingPonger(dw, ppimg, 390, 180, 2)
	createPingPonger(dw, ppimg, 400, 180, 3)
}

func createCharacter(dw primen.World, ebimg *ebiten.Image) {
	sn := primen.NewRootAnimatedSpriteNode(dw, primen.Layer0, 14, &graphics.TiledAnimation{
		Clips: []graphics.TiledAnimationClip{
			{
				Image: ebimg,
				Name:  "default",
				Frames: []image.Rectangle{
					image.Rect(0, 0, 80, 110),      // 0
					image.Rect(80, 0, 80*2, 110),   // 1
					image.Rect(80*2, 0, 80*3, 110), // 2
					image.Rect(80*3, 0, 80*4, 110), // 3
					image.Rect(0, 0, 80, 110),      // 0
				},
				ClipMode: graphics.AnimLoop,
			},
		},
	})
	s := sn.Transform()
	s.SetX(300).SetY(200)
	sn.SpriteAnim().PlayClip("default")
}

func createPingPonger(dw primen.World, ebimg *ebiten.Image, x, y float64, frame int) {
	sn := primen.NewRootAnimatedSpriteNode(dw, primen.Layer0, 24, &graphics.TiledAnimation{
		Clips: []graphics.TiledAnimationClip{
			{
				Name:  "default",
				Image: ebimg,
				Frames: []image.Rectangle{
					image.Rect(8*0, 0, 8*1, 32),
					image.Rect(8*1, 0, 8*2, 32),
					image.Rect(8*2, 0, 8*3, 32),
					image.Rect(8*3, 0, 8*4, 32),
				},
				ClipMode: graphics.AnimPingPong,
			},
		},
	})
	s := sn.Transform()
	s.SetX(x).SetY(y).SetAngle(math.Pi / 4)
	sn.SpriteAnim().PlayClipFrame("default", frame)
}
