package main

import (
	"context"
	"math"

	"github.com/gabstv/tau"
	"github.com/gabstv/tau/core"
	"github.com/gabstv/tau/examples/layers/res"
	"github.com/gabstv/tau/io"
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
		atlas.Get("box1"),
		atlas.Get("box2"),
		atlas.Get("box3"),
		atlas.Get("box4"),
	}
	core.DebugDraw = true
	//
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	//
	engine := tau.NewEngine(&tau.NewEngineInput{
		Width:           640,
		Height:          480,
		FS:              fs,
		Title:           "New Transform Test",
		Scale:           2,
		FixedResolution: true,
		Resizable:       true,
		OnReady: func(e *tau.Engine) {
			dogamesetup(ctx, e, spbgs)
		},
	})
	if err := engine.Run(); err != nil {
		println(err.Error())
	}
}

func dogamesetup(ctx context.Context, engine *tau.Engine, bgs []*ebiten.Image) {
	select {
	case <-ctx.Done():
		return
	case <-engine.Ready():
	}
	spr := tau.NewSprite(engine.Default(), bgs[0], tau.Layer0, nil)
	spr.Transform.X = 100
	spr.Transform.Y = 100
	spr.Transform.ScaleX = .5
	spr.TauSprite.OriginX = .5
	spr.TauSprite.OriginY = .5
	engine.Default().AddComponentToEntity(spr.Entity, engine.Default().Component(core.CNRotation), &core.Rotation{
		Speed: math.Pi / 16,
	})
	spr2 := tau.NewSprite(engine.Default(), bgs[1], tau.Layer0, spr.Transform)
	spr2.Transform.X = 10
	spr2.Transform.Y = 7

	spr3 := tau.NewSprite(engine.Default(), bgs[2], tau.Layer0, spr2.Transform)
	spr3.Transform.X = 16
	spr3.Transform.Y = 16
	spr3.Transform.ScaleX = 2
}
