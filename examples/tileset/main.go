// +build example

package main

import (
	"context"
	"fmt"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func main() {
	ebiten.SetRunnableInBackground(true)
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
		OnReady: func(e *primen.Engine) {
			dogamesetup(ctx, e, spbgs)
		},
	})
	if err := engine.Run(); err != nil {
		println(err.Error())
	}
}

func dogamesetup(ctx context.Context, engine *primen.Engine, bgs []*ebiten.Image) {
	select {
	case <-ctx.Done():
		return
	case <-engine.Ready():
	}
	tset := primen.NewTileSet(engine.Root(nil), primen.Layer0)
	tset.SetDB(bgs)
	tset.SetCellSize(32, 32)
	tset.SetTilesYX([][]int{
		{0, 0, 0, 0, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3},
		{2, 1, 3, 0, 1, 2, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2},
		{0, 0, 2, 3, 1, 0, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1},
		{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0},
		{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 3},
		{1, 3, 1, 1, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 1, 3, 2},
		{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 1},
		{0, 0, 0, 0, 0, 3, 0, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 0},
		{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 1, 2, 1, 2, 3, 1, 2, 3},
		{0, 2, 0, 0, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2},
		{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 1},
		{0, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 0},
		{0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 2, 2, 1, 2, 3, 1, 2, 3},
		{1, 2, 1, 1, 3, 3, 2, 3, 3, 1, 2, 2, 3, 2, 1, 2, 3, 1, 2, 2},
		{0, 2, 1, 0, 3, 3, 2, 3, 3, 1, 2, 3, 3, 2, 1, 2, 3, 1, 2, 2},
	})
	tset.SetOrigin(.5, .5)
	tset.SetPos(float64(engine.Width())/2, float64(engine.Height())/2)
	dd := primen.NewTransform(engine.Root(nil))
	help := false
	helpoff := `press 'h' for help`
	helpon := `press 'h' to hide
  keys:
    [1] 50% scale; [2] 100% scale, [3] 200% scale
    [-] decrease scale; [+] increase scale
    WASD, Arrow Keys: move tileset
    [q] [e] rotate`
	dd.UpsertDrawFns(nil, nil, func(ctx core.Context, e ecs.Entity) {
		fps := ebiten.CurrentFPS()
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), fmt.Sprintf("%.2f fps", fps), 12, 12)
		if !help {
			ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), helpoff, 12, 24)
		} else {
			ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), helpon, 12, 24)
		}
	})
	dd.UpsertFns(func(ctx core.Context, e ecs.Entity) {
		//tset.SetPos(float64(engine.Width())/2, float64(engine.Height())/2)
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			tset.SetScale2(.5)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			tset.SetScale2(1)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			tset.SetScale2(2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			tset.SetAngle(tset.Angle() + ctx.DT()*-2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			tset.SetAngle(tset.Angle() + ctx.DT()*2)
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
			tset.SetY(tset.Y() - ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
			tset.SetY(tset.Y() + ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
			tset.SetX(tset.X() - ctx.DT()*100)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
			tset.SetX(tset.X() + ctx.DT()*100)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
			s, _ := tset.Scale()
			tset.SetScale2(s + 0.25)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
			s, _ := tset.Scale()
			tset.SetScale2(s - 0.25)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyH) {
			help = !help
		}
	}, nil, nil)
}
