// +build example

package main

import (
	"context"
	"math"
	"math/rand"
	"unsafe"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/examples/layers/layerexample"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	// "github.com/pkg/profile" // enable line 24
)

func main() {
	ebiten.SetRunnableInBackground(true)
	//defer profile.Start(profile.MemProfile).Stop()
	fs := res.FS()
	container := io.NewContainer(context.TODO(), fs)
	<-container.Load("public/atlas.dat")
	atlas, err := container.GetAtlas("public/atlas.dat")
	if err != nil {
		panic(err)
	}
	spbgs := []*ebiten.Image{
		atlas.GetSubImage("box1").Image,
		atlas.GetSubImage("box2").Image,
		atlas.GetSubImage("box3").Image,
		atlas.GetSubImage("box4").Image,
	}
	spfgs := []*ebiten.Image{
		atlas.GetSubImage("l1").Image,
		atlas.GetSubImage("l2").Image,
		atlas.GetSubImage("l3").Image,
		atlas.GetSubImage("l4").Image,
	}
	//
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	//
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:         640,
		Height:        480,
		FS:            fs,
		Title:         "Layers Test",
		Scale:         0.5,
		Resizable:     true,
		MaxResolution: true,
		OnReady: func(e primen.Engine) {
			dogamesetup(ctx, e, spbgs, spfgs)
		},
	})
	engine.SetDebugTPS(true)
	engine.AddEventListener("act_of_nature", func(eventName string, e core.Event) {
		println("act of nature happened!")
		println(e.Data.(ecs.Entity))
		// if globalScale != 1 {
		// 	globalScale = 1
		// } else {
		// 	globalScale = -2
		// }
	})
	if err := engine.Run(); err != nil {
		println(err.Error())
	}
}

func dogamesetup(ctx context.Context, engine primen.Engine, bgs, fgs []*ebiten.Image) {
	select {
	case <-ctx.Done():
		return
	case <-engine.Ready():
	}

	w := engine.NewWorldWithDefaults(0)

	layerexample.GetOrbitalMovementSystem(w).SetBgs(bgs)
	layerexample.GetOrbitalMovementSystem(w).SetFgs(fgs)

	rand.Seed(112358)

	rootnode := primen.NewRootFnNode(w)
	rootnode.Transform().SetX(320 / 2).SetY(240 / 2)

	rootnode.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		ttx := components.GetTransformComponentData(w, e)
		xx := uintptr(unsafe.Pointer(ttx))
		_ = xx
		ttx.SetX(float64(ctx.Engine().Width() / 2)).SetY(float64(ctx.Engine().Height() / 2))
	}

	nrings := 5
	nitems := 50

	w.Listen(ecs.EvtComponentsResized, func(e ecs.Event) {
		println("COMPONENT RESIZED", e.ComponentID, e.ComponentName)
	})

	for ring := 0; ring < nrings; ring++ {
		for itemi := 0; itemi < nitems; itemi++ {
			layer := rand.Intn(4)
			bgnode := primen.NewChildSpriteNode(rootnode, primen.Layer(layer))
			if rand.Float64() < 0.1 {
				bgnode.Sprite().RotateHue(rand.Float64() * (math.Pi * 2))
			}
			if rand.Float64() < 0.06 {
				println("COMPOSITE LIGHTER")
				bgnode.Sprite().SetCompositeMode(ebiten.CompositeModeLighter)
			} else if rand.Float64() < 0.03 {
				println("COMPOSITE XOR")
				bgnode.Sprite().SetCompositeMode(ebiten.CompositeModeXor)
			}
			bgnode.Sprite().SetOrigin(.5, .5).SetImage(bgs[layer])
			fgnode := primen.NewChildSpriteNode(bgnode, primen.Layer(layer))
			fgnode.Sprite().SetOrigin(.5, .5).SetImage(fgs[layer])
			layerexample.SetOrbitalMovementComponentData(w, bgnode.Entity(), layerexample.OrbitalMovement{
				Dx:          float64(itemi+1)*30 + rand.Float64()*10,
				Dy:          float64(itemi+1)*30 + rand.Float64()*10,
				ChildSprite: fgnode.Entity(),
				R:           math.Pi * rand.Float64() * 2,
				Speed:       float64(5-itemi)/4 + rand.Float64()/4,
				Ox:          (rand.Float64() - 0.5) * 5,
				Oy:          (rand.Float64() - 0.5) * 5,
				AngleR:      rand.Float64(),
			})
			if rand.Float64() < 0.12 {
				println("HUE SHIFTER")
				layerexample.GetOrbitalMovementComponentData(w, bgnode.Entity()).HueShift = true
			}
		}
	}

	rootnode = nil

	fnss := primen.NewRootFnNode(w)

	fnss.Function().Draw = func(ctx core.DrawCtx, e ecs.Entity) {
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "d: toggle debug draw", 10, 20)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "p: toggle pause", 10, 20+12*1)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "n: next frame (while paused)", 10, 20+12*2)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "[: decrease radius", 10, 20+12*3)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "]: increase radius", 10, 20+12*4)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "s: high dpi", 10, 20+12*5)
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), "a: low dpi", 10, 20+12*6)
	}

	fnss.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		mvsys := layerexample.GetOrbitalMovementSystem(w)
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			debug.Draw = !debug.Draw
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			mvsys.TogglePause()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyN) && mvsys.Paused() {
			mvsys.PushFrame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRightBracket) {
			mvsys.AddScale()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeftBracket) {
			mvsys.SubScale()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			engine.SetScreenScale(ebiten.DeviceScaleFactor())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			engine.SetScreenScale(.5)
		}
	}
}
