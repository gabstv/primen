package main

import (
	"context"
	"math"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
)

func main() {
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
	core.DebugDraw = true
	//
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	//
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:           640,
		Height:          480,
		FS:              fs,
		Title:           "New Transform Test",
		Scale:           0.5,
		FixedResolution: true,
		Resizable:       true,
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
	spr := primen.NewSprite(engine.Root(nil), bgs[0], primen.Layer0)
	spr.SetX(100)
	spr.SetY(100)
	spr.SetScale(.5, 1)
	spr.SetOrigin(.5, .5)
	engine.Default().AddComponentToEntity(spr.Entity(), engine.Default().Component(core.CNRotation), &core.Rotation{
		Speed: math.Pi / 16,
	})
	spr2 := primen.NewSprite(spr, bgs[1], primen.Layer0)
	spr2.SetPos(10, 7)

	spr3 := primen.NewSprite(spr2, bgs[2], primen.Layer0)
	spr3.SetPos(16, 16)
	spr3.SetScale(2, 1)
}
