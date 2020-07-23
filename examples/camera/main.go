package main

import (
	"bytes"
	"image"
	_ "image/png"
	"io/ioutil"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
)

var (
	bgimage *ebiten.Image
	fgimage *ebiten.Image
)

func main() {
	load()
	debug.Draw = true
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		OnReady: func(e primen.Engine) {
			w := e.NewWorldWithDefaults(0)
			wtrn := primen.NewRootNode(w)
			ctr := primen.NewChildNode(wtrn)
			components.SetCameraComponentData(w, ctr.Entity(), components.NewCamera(e.NewScreenOffsetDrawTarget(core.DrawMaskDefault)))
			cam := components.GetCameraComponentData(w, ctr.Entity())
			cam.SetViewRect(e.SizeVec())
			components.SetFollowTransformComponentData(w, ctr.Entity(), components.FollowTransform{})
			c := components.GetFollowTransformComponentData(w, ctr.Entity())
			//c.SetDeadZone(geom.Rect{geom.ZV, geom.Vec{100, 100}}.At(e.SizeVec().Scaled(.5).Sub(geom.Vec{50, 50})))
			tiled := primen.NewChildTileSetNode(wtrn, primen.Layer0, []*ebiten.Image{bgimage}, 16, 16, 256, 256, make([]int, 16*16))
			_ = tiled
			dude := primen.NewChildSpriteNode(wtrn, primen.Layer1)
			dude.Sprite().SetImage(fgimage).SetOrigin(.5, .5)
			dude.Transform().SetPos(e.SizeVec().Scaled(.5))
			ctr.Transform().SetPos(e.SizeVec().Scaled(.5))
			c.SetTarget(dude.Entity())
			c.SetBounds(geom.Rect{
				Min: e.SizeVec().Scaled(.5),
				Max: geom.Vec{2048, 2048},
			})
			c.SetDeadZone(geom.Vec{100, 40})
			//
			fnn := primen.NewRootFnNode(w)
			fnn.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
				dtr := dude.Transform()
				if ebiten.IsKeyPressed(ebiten.KeyRight) {
					dtr.SetX(dtr.X() + (ctx.DT() * 200))
				}
				if ebiten.IsKeyPressed(ebiten.KeyLeft) {
					dtr.SetX(dtr.X() - (ctx.DT() * 200))
				}
				if ebiten.IsKeyPressed(ebiten.KeyDown) {
					dtr.SetY(dtr.Y() + (ctx.DT() * 200))
				}
				if ebiten.IsKeyPressed(ebiten.KeyUp) {
					dtr.SetY(dtr.Y() - (ctx.DT() * 200))
				}
			}
		},
	})
	engine.Run()
}

func load() {
	b, err := ioutil.ReadFile("cambg.png")
	if err != nil {
		panic(err)
	}
	rawbgimage, _, _ := image.Decode(bytes.NewReader(b))
	bgimage, _ = ebiten.NewImageFromImage(rawbgimage, ebiten.FilterNearest)
	b2, err := ioutil.ReadFile("target.png")
	if err != nil {
		panic(err)
	}
	rawfgimage, _, _ := image.Decode(bytes.NewReader(b2))
	fgimage, _ = ebiten.NewImageFromImage(rawfgimage, ebiten.FilterNearest)
}
