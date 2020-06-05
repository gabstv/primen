package main

// https://www.kenney.nl/assets/platformer-characters-1

import (
	"fmt"
	"image"
	_ "image/png"
	"math"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("zombie_kenney.png", ebiten.FilterDefault)
	ppimg, _, _ := ebitenutil.NewImageFromFile("ping_pong.png", ebiten.FilterDefault)

	engine := primen.NewEngine(&primen.NewEngineInput{
		Title:  "Basic Animation",
		Width:  640,
		Height: 480,
		Scale:  1,
	})

	core.DebugDraw = true

	dw := engine.Default()
	sc := dw.Component(core.CNDrawable)
	ac := dw.Component(core.CNSpriteAnimation)
	createCharacter(dw, sc, ac, ebimg)
	createPingPonger(dw, sc, ac, ppimg)

	s0 := dw.NewSystem("", 0, func(ctx ecs.Context) {
		screen := ctx.World().Get("screen").(*ebiten.Image)
		fps := ebiten.CurrentFPS()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f fps", fps), 0, 0)
	})
	s0.AddTag(primen.WorldTagDraw)

	engine.Run()
}

func createCharacter(dw *ecs.World, spriteComp *ecs.Component, animComp *ecs.Component, ebimg *ebiten.Image) {
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, spriteComp, &core.Sprite{
		Image:  ebimg.SubImage(image.Rect(0, 0, 80, 110)).(*ebiten.Image),
		X:      300,
		Y:      200,
		Angle:  0,
		ScaleX: 1,
		ScaleY: 1,
	})
	dw.AddComponentToEntity(e, animComp, &core.SpriteAnimation{
		Enabled: true,
		Playing: true,
		Anim: &core.TiledAnimation{
			Clips: []core.TiledAnimationClip{
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
					ClipMode: core.AnimLoop,
				},
			},
		},
		Fps: 24,
	})
}

func createPingPonger(dw *ecs.World, spriteComp *ecs.Component, animComp *ecs.Component, ebimg *ebiten.Image) {
	e := dw.NewEntity()
	dw.AddComponentToEntity(e, spriteComp, &core.Sprite{
		Image:  ebimg.SubImage(image.Rect(0, 0, 8, 32)).(*ebiten.Image),
		X:      370,
		Y:      180,
		Angle:  math.Pi / 4,
		ScaleX: 1,
		ScaleY: 1,
	})
	dw.AddComponentToEntity(e, animComp, &core.SpriteAnimation{
		Enabled: true,
		Playing: true,
		Anim: &core.TiledAnimation{
			Clips: []core.TiledAnimationClip{
				{
					Name:  "default",
					Image: ebimg,
					Frames: []image.Rectangle{
						image.Rect(8*0, 0, 8*1, 32),
						image.Rect(8*1, 0, 8*2, 32),
						image.Rect(8*2, 0, 8*3, 32),
						image.Rect(8*3, 0, 8*4, 32),
					},
					ClipMode: core.AnimPingPong,
				},
			},
		},
		Fps: 24,
	})
}
