package main

import (
	"bytes"
	"image"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
	//"github.com/pkg/profile"
)

var (
	bgimage *ebiten.Image
	fgimage *ebiten.Image
	arimage *ebiten.Image
)

func main() {
	double := false
	if len(os.Args) > 1 && os.Args[1] == "split" {
		double = true
	}
	//defer profile.Start(profile.CPUProfile).Stop()
	load()
	//debug.Draw = true
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		OnReady: func(e primen.Engine) {
			if !double {
				singleCameraExample(e)
			} else {
				doubleCameraExample(e)
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
	b3, err := ioutil.ReadFile("arrow.png")
	if err != nil {
		panic(err)
	}
	rawarrow, _, _ := image.Decode(bytes.NewReader(b3))
	arimage, _ = ebiten.NewImageFromImage(rawarrow, ebiten.FilterNearest)
}

func singleCameraExample(e primen.Engine) {
	e.SetDebugTPS(true)
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
	ar1 := primen.NewChildSpriteNode(dude, primen.Layer2)
	ar1.Transform().SetScale(3, 3)
	ar1.Sprite().SetImage(arimage).SetOrigin(.5, .5)
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
}

func doubleCameraExample(e primen.Engine) {
	e.SetDebugTPS(true)
	t1id := e.NewDrawTarget(core.DrawMaskDefault, geom.Rect{Max: geom.Vec{.495, 1}}, ebiten.FilterDefault)
	t2id := e.NewDrawTarget(core.DrawMaskDefault, geom.Rect{Min: geom.Vec{0.505, 0}, Max: geom.Vec{1, 1}}, ebiten.FilterDefault)

	w := e.NewWorldWithDefaults(0)
	wtrn := primen.NewRootNode(w)

	tiled := primen.NewChildTileSetNode(wtrn, primen.Layer0, []*ebiten.Image{bgimage}, 16, 16, 256, 256, make([]int, 16*16))
	_ = tiled

	var player1 *primen.SpriteNode
	var player2 *primen.SpriteNode
	var ar1, ar2 *primen.SpriteNode
	// player 1
	{
		player1 = primen.NewChildSpriteNode(wtrn, primen.Layer1)
		player1.Sprite().SetImage(fgimage).SetOrigin(.5, .5)
		player1.Transform().SetPos(e.SizeVec().Scaled(.5))
		ar1 = primen.NewChildSpriteNode(player1, primen.Layer2)
		ar1.Transform().SetScale(3, 3)
		ar1.Sprite().SetImage(arimage).SetOrigin(.5, .5)
	}
	// player 2
	{
		player2 = primen.NewChildSpriteNode(wtrn, primen.Layer1)
		player2.Sprite().SetImage(fgimage).SetOrigin(.5, .5).RotateHue(math.Pi)
		player2.Transform().SetPos(e.SizeVec().Scaled(.5).Add(geom.Vec{0, 50}))
		ar2 = primen.NewChildSpriteNode(player2, primen.Layer2)
		ar2.Transform().SetScale(3, 3)
		ar2.Sprite().SetImage(arimage).SetOrigin(.5, .5)
	}
	// camera 1
	{
		camtr := primen.NewChildNode(wtrn)
		components.SetCameraComponentData(w, camtr.Entity(), components.NewCamera(t1id))
		cam := components.GetCameraComponentData(w, camtr.Entity())
		cam.SetViewRect(e.SizeVec().ScaledXY(.5, 1))
		components.SetFollowTransformComponentData(w, camtr.Entity(), components.FollowTransform{})
		ftr := components.GetFollowTransformComponentData(w, camtr.Entity())
		ftr.SetTarget(player1.Entity())
		ftr.SetBounds(geom.Rect{
			Min: e.SizeVec().Scaled(.5),
			Max: geom.Vec{2048, 2048},
		})
		ftr.SetDeadZone(geom.Vec{80, 40})
		camtr.Transform().SetPos(player1.Transform().Pos())
	}
	// camera 2
	{
		camtr := primen.NewChildNode(wtrn)
		components.SetCameraComponentData(w, camtr.Entity(), components.NewCamera(t2id))
		cam := components.GetCameraComponentData(w, camtr.Entity())
		cam.SetViewRect(e.SizeVec().ScaledXY(.5, 1))
		components.SetFollowTransformComponentData(w, camtr.Entity(), components.FollowTransform{})
		ftr := components.GetFollowTransformComponentData(w, camtr.Entity())
		ftr.SetTarget(player2.Entity())
		ftr.SetBounds(geom.Rect{
			Min: e.SizeVec().Scaled(.5),
			Max: geom.Vec{2048, 2048},
		})
		ftr.SetDeadZone(geom.Vec{80, 40})
		camtr.Transform().SetPos(player2.Transform().Pos())
	}
	//
	fnn := primen.NewRootFnNode(w)
	fnn.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		a2a := player1.Transform().Pos().Sub(player2.Transform().Pos()).Angle()
		a1a := player2.Transform().Pos().Sub(player1.Transform().Pos()).Angle()
		ar1.Transform().SetAngle(a1a)
		ar2.Transform().SetAngle(a2a)
		dtr := player2.Transform()
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
		dtr2 := player1.Transform()
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			dtr2.SetX(dtr2.X() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			dtr2.SetX(dtr2.X() - (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			dtr2.SetY(dtr2.Y() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			dtr2.SetY(dtr2.Y() - (ctx.DT() * 200))
		}
	}
}
