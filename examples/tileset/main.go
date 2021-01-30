// +build example

package main

import (
	"context"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	ebiten.SetRunnableOnUnfocused(true)
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
	//
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	//
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:         640,
		Height:        480,
		FS:            fs,
		Title:         "TileSet Test",
		Scale:         1,
		Resizable:     true,
		MaxResolution: true,
		OnReady: func(e primen.Engine) {
			dogamesetup(ctx, e, spbgs)
		},
	})
	engine.SetDebugTPS(true)
	if err := engine.Run(); err != nil {
		println(err.Error())
	}
}

func dogamesetup(ctx context.Context, engine primen.Engine, bgs []*ebiten.Image) {
	select {
	case <-ctx.Done():
		return
	case <-engine.Ready():
	}
	w := engine.NewWorldWithDefaults(0)
	tset := primen.NewRootTileSetNode(w, primen.Layer0, bgs, 15, 20, 32, 32, nil)
	tset.TileSet().SetCells([]int{
		0, 0, 0, 0, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3,
		2, 1, 3, 0, 1, 2, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2,
		0, 0, 2, 3, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1,
		0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0,
		0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3,
		1, 3, 1, 1, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 1, 3, 2,
		0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 1,
		0, 0, 0, 0, 0, 3, 0, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 0,
		1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 3,
		0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2,
		1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1,
		0, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0,
		0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 2, 2, 1, 2, 3, 1, 2, 3,
		1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2,
		0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2,
	})
	// tset := primen.NewTileSet(engine.Root(nil), primen.Layer0)
	// tset.SetDB(bgs)
	// tset.SetCellSize(32, 32)
	// tset.SetTilesYX([][]int{
	// 	{0, 0, 0, 0, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3},
	// 	{2, 1, 3, 0, 1, 2, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2},
	// 	{0, 0, 2, 3, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1},
	// 	{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0},
	// 	{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3},
	// 	{1, 3, 1, 1, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 1, 3, 2},
	// 	{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 1},
	// 	{0, 0, 0, 0, 0, 3, 0, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 0},
	// 	{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 3},
	// 	{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2},
	// 	{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1},
	// 	{0, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0},
	// 	{0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 2, 2, 1, 2, 3, 1, 2, 3},
	// 	{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2},
	// 	{0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2},
	// })
	tset.TileSet().SetOrigin(.5, .5)
	tset.Transform().SetX(float64(engine.Width()) / 2).SetY(float64(engine.Height()) / 2)
	//tset.SetOrigin(.5, .5)
	//tset.SetPos(float64(engine.Width())/2, float64(engine.Height())/2)
	//dd := primen.NewTransform(engine.Root(nil))
	dd := primen.NewRootFnNode(w)
	help := false
	helpoff := `press 'h' for help`
	helpon := `press 'h' to hide
  keys:
    [1] 50% scale; [2] 100% scale, [3] 200% scale
    [-] decrease scale; [+] increase scale
    WASD, Arrow Keys: move tileset
	[q] [e] rotate`
	dd.Function().Draw = func(ctx core.DrawCtx, e ecs.Entity) {
		if !help {
			ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), helpoff, 10, 24)
		} else {
			ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), helpon, 10, 24)
		}
	}
	dd.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			tset.Transform().SetScale(.5, .5)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			tset.Transform().SetScale(1, 1)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			tset.Transform().SetScale(2, 2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			a := tset.Transform().Angle()
			tset.Transform().SetAngle(a + ctx.DT()*-2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			a := tset.Transform().Angle()
			tset.Transform().SetAngle(a + ctx.DT()*2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
			t := tset.Transform()
			t.SetY(t.Y() - ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
			t := tset.Transform()
			t.SetY(t.Y() + ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
			t := tset.Transform()
			t.SetX(t.X() - ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
			t := tset.Transform()
			t.SetX(t.X() + ctx.DT()*100)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
			t := tset.Transform()
			sx := t.ScaleX()
			t.SetScale(sx+.25, sx+.25)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
			t := tset.Transform()
			sx := t.ScaleX()
			t.SetScale(sx-.25, sx-.25)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyH) {
			help = !help
		}
	}
}
