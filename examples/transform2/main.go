package main

import (
	"context"
	"math"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/debug"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/gabstv/ebiten"
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
	debug.Draw = true
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
		OnReady: func(e primen.Engine) {
			dogamesetup(ctx, e, spbgs)
		},
	})
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
	nspr := primen.NewRootSpriteNode(w, primen.Layer0)
	tr := nspr.Transform()
	spr := nspr.Sprite()
	spr.SetImage(bgs[0]).SetOrigin(.5, .5)
	tr.SetX(100).SetY(100).SetScale(.5, 1)

	// engine.Default().AddComponentToEntity(spr.Entity(), engine.Default().Component(core.CNRotation), &core.Rotation{
	// 	Speed: math.Pi / 16,
	// })

	{
		// spr2
		nspr2 := primen.NewChildSpriteNode(nspr, primen.Layer0)
		nspr2.Sprite().SetImage(bgs[1])
		nspr2.Transform().SetX(10).SetY(7)

		// spr3
		nspr3 := primen.NewChildSpriteNode(nspr2, primen.Layer0)
		nspr3.Sprite().SetImage(bgs[2])
		nspr3.Transform().SetX(16).SetY(16).SetScale(2, 1)

		components.SetFunctionComponentData(w, nspr3.Entity(), components.Function{
			Update: func(ctx core.UpdateCtx, e ecs.Entity) {
				dd := components.GetTransformComponentData(w, e)
				dd.SetAngle(dd.Angle() + .1*math.Pi*ctx.DT())
			},
		})

		// spr4
		nspr4 := primen.NewChildSpriteNode(nspr3, primen.Layer0)
		nspr4.Sprite().SetImage(bgs[0])
		nspr4.Transform().SetX(32).SetY(32).SetScale(.25, .5)
	}

	components.SetFunctionComponentData(w, nspr.Entity(), components.Function{
		Update: func(ctx core.UpdateCtx, e ecs.Entity) {
			dd := components.GetTransformComponentData(w, e)
			dd.SetAngle(dd.Angle() - .5*math.Pi*ctx.DT())
		},
	})
}
