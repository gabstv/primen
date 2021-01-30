package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/easing"
	"github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var pimg *ebiten.Image
var limg *ebiten.Image

func main() {
	pimg, _, _ = ebitenutil.NewImageFromFile("../shared/particle.png")
	limg, _, _ = ebitenutil.NewImageFromFile("../shared/primen_logo@0.5x.png")

	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		Scale:  ebiten.DeviceScaleFactor(),
		Title:  "Tweening",
		OnReady: func(e primen.Engine) {
			w := e.NewWorldWithDefaults(0)
			root := primen.NewRootFnNode(w)
			dsf := ebiten.DeviceScaleFactor()
			root.Transform().SetScale(dsf*2, dsf*2)
			root.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
				root.Transform().SetX(float64(ctx.Engine().Width() / 2))
				root.Transform().SetY(float64(ctx.Engine().Height() / 2))
			}
			spr := primen.NewChildSpriteNode(root, primen.Layer0)
			spr.Transform().SetY(-70)
			spr.TrTweening().SetTween("ltr", components.TrTweenX, -100, 100, 2, easing.OutQuart)
			spr.TrTweening().SetTween("rtl", components.TrTweenX, 100, -100, 3, easing.OutElasticFunction(.5))
			spr.TrTweening().SetDoneCallback(func(name string) {
				if name == "ltr" {
					spr.TrTweening().Play("rtl")
				} else if name == "rtl" {
					spr.TrTweening().Play("ltr")
				}
			})
			spr.TrTweening().Play("ltr")
			spr.Sprite().SetImage(pimg).SetOrigin(.5, .5)
			//
			spr2 := primen.NewChildSpriteNode(root, primen.Layer0)
			spr2.Transform().SetY(70)
			spr2.TrTweening().SetTween("ltr", components.TrTweenX, -120, 120, 1, easing.OutBounce)
			spr2.TrTweening().SetTween("rtl", components.TrTweenX, 120, -120, 1, easing.InOutCubic)
			spr2.TrTweening().SetDoneCallback(func(name string) {
				if name == "ltr" {
					spr2.TrTweening().Play("rtl")
				} else if name == "rtl" {
					spr2.TrTweening().Play("ltr")
				}
			})
			spr2.TrTweening().Play("ltr")
			spr2.Sprite().SetImage(pimg).SetOrigin(.5, .5)
			//
			spr3 := primen.NewChildSpriteNode(root, primen.Layer0)
			spr3.Transform().SetScale(1/(dsf*6), 1/(dsf*6))
			spr3.TrTweening().SetTween("loop", components.TrTweenRotation, 0, math.Pi*6, 3, easing.InOutElastic)
			spr3.TrTweening().SetTween("loop2", components.TrTweenRotation, 0, math.Pi*2, 2, easing.InOutSine)
			spr3.TrTweening().SetDoneCallback(func(name string) {
				if name == "scaleless" {
					spr3.TrTweening().SetTween("scaleplus", components.TrTweenScaleXY, spr3.Transform().ScaleX(), spr3.Transform().ScaleX()+.1, 1, easing.OutBounce)
					spr3.TrTweening().Play("scaleless")
				}
				go func() {
					time.Sleep(time.Second)
					e.RunFn(func() {
						if rand.Intn(2) == 0 {
							spr3.TrTweening().Play("loop")
						} else {
							spr3.TrTweening().Play("loop2")
						}
					})
				}()
			})
			spr3.TrTweening().Play("loop")
			spr3.Sprite().SetImage(limg).SetOrigin(.5, .5)
			////////
			spr4 := primen.NewChildSpriteNode(root, primen.Layer0)
			spr4.Transform().SetScale(1/(dsf*11), 1/(dsf*11)).SetY(100)
			spr4.TrTweening().SetTween("tobig", components.TrTweenScaleXY, 1/(dsf*14), 1/(dsf*11), 1, easing.OutBounce)
			spr4.TrTweening().SetTween("tosmall", components.TrTweenScaleXY, 1/(dsf*11), 1/(dsf*14), 1, easing.OutBounce)
			spr4.TrTweening().SetDoneCallback(func(name string) {
				if name == "tobig" {
					spr4.TrTweening().Play("tosmall")
				} else {
					spr4.TrTweening().Play("tobig")
				}
			})
			spr4.TrTweening().Play("tosmall")
			spr4.Sprite().SetImage(limg).SetOrigin(.5, .5)
		},
	})
	engine.Run()
}
