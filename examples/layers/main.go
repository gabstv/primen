package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
	"github.com/gabstv/tau/core"
	"github.com/gabstv/tau/examples/layers/res"
	"github.com/gabstv/tau/io"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var movementPaused = false
var xframes = make(chan struct{}, 30)

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
	spfgs := []*ebiten.Image{
		atlas.Get("l1"),
		atlas.Get("l2"),
		atlas.Get("l3"),
		atlas.Get("l4"),
	}
	//
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	//
	engine := tau.NewEngine(&tau.NewEngineInput{
		Width:         640,
		Height:        480,
		FS:            fs,
		Title:         "Layers Test",
		Scale:         2,
		Resizable:     true,
		MaxResolution: true,
		OnReady: func(e *tau.Engine) {
			dogamesetup(ctx, e, spbgs, spfgs)
		},
	})
	if err := engine.Run(); err != nil {
		println(err.Error())
	}
}

type OrbitalMovement struct {
	Speed       float64
	Dx          float64
	Dy          float64
	Ox          float64
	Oy          float64
	R           float64
	AngleR      float64
	ChildSprite *tau.Sprite
}

func dogamesetup(ctx context.Context, engine *tau.Engine, bgs, fgs []*ebiten.Image) {
	select {
	case <-ctx.Done():
		return
	case <-engine.Ready():
	}

	movecompname := "__movement_comp"

	movecs := &core.BasicCS{
		SysName: "__movement_system",
		GetComponents: func(w *ecs.World) []*ecs.Component {
			return []*ecs.Component{
				core.UpsertComponent(w, ecs.NewComponentInput{
					Name: movecompname,
				}),
				w.Component(core.CNTransform),
				w.Component(core.CNDrawLayer),
				w.Component(core.CNDrawable),
			}
		},
		SysPriority: -3,
		SysExec: func(ctx core.Context) {
			if movementPaused {
				select {
				case <-xframes:
					// exec a frame
				default:
					return
				}
			}
			trc := ctx.World().Component(core.CNTransform)
			dlc := ctx.World().Component(core.CNDrawLayer)
			moc := ctx.World().Component(movecompname)
			spc := ctx.World().Component(core.CNDrawable)
			dt := ctx.DT()
			//
			for _, match := range ctx.System().View().Matches() {
				sprite := match.Components[spc].(*core.Sprite)
				transform := match.Components[trc].(*core.Transform)
				drawlayer := match.Components[dlc].(*core.DrawLayer)
				movecomp := match.Components[moc].(*OrbitalMovement)
				movecomp.R += movecomp.Speed * dt
				xx := math.Cos(movecomp.R) * movecomp.Dx
				yy := math.Sin(movecomp.R) * movecomp.Dy
				transform.X = movecomp.Ox + xx
				transform.Y = movecomp.Oy + yy
				transform.Angle += dt * (math.Pi / 4) * movecomp.AngleR
				if rand.Float64() < 0.001 {
					newlayer := rand.Intn(4)
					drawlayer.Layer = core.LayerIndex(newlayer)
					sprite.Image = bgs[newlayer]
					sprite.Bounds = sprite.Image.Bounds()
					//drawlayer.ZIndex = 1
					movecomp.ChildSprite.TauSprite.Image = fgs[newlayer]
					movecomp.ChildSprite.TauSprite.Bounds = fgs[newlayer].Bounds()
					movecomp.ChildSprite.DrawLayer.Layer = core.LayerIndex(newlayer)
					//movecomp.ChildSprite.OriginX
				}
			}
		},
	}
	//
	_ = movecs.Components(engine.Default())
	core.SetupSystem(engine.Default(), movecs)

	rand.Seed(112358)

	root := tau.NewTransform(engine.Default(), nil)
	root.TauTransform.X = 320 / 2
	root.TauTransform.Y = 240 / 2

	root.UpsertFns(func(ctx core.Context, e ecs.Entity) {
		t := ctx.World().Component(core.CNTransform).Data(e).(*core.Transform)
		t.X = float64(ctx.Engine().Width() / 2)
		t.Y = float64(ctx.Engine().Height() / 2)
	}, nil, nil)

	for i := 0; i < 4; i++ {
		for j := 0; j < 20; j++ {
			//ri := rand.Intn(4)
			rl := rand.Intn(4)
			bgs := tau.NewSprite(engine.Default(), bgs[rl], core.LayerIndex(rl), root.TauTransform)
			bgs.TauSprite.OriginX = .5
			bgs.TauSprite.OriginY = .5
			fgs := tau.NewSprite(engine.Default(), fgs[rl], core.LayerIndex(rl), bgs.Transform)
			fgs.TauSprite.OriginX = .5
			fgs.TauSprite.OriginY = .5
			//fgs.Transform.Angle = -math.Pi * 0.5
			mvc := &OrbitalMovement{
				Dx:          float64(i+1)*30 + rand.Float64()*10,
				Dy:          float64(i+1)*30 + rand.Float64()*10,
				ChildSprite: fgs,
				R:           math.Pi * rand.Float64() * 2,
				Speed:       float64(5-i)/4 + rand.Float64()/4,
				Ox:          (rand.Float64() - 0.5) * 5,
				Oy:          (rand.Float64() - 0.5) * 5,
				AngleR:      rand.Float64(),
			}
			engine.Default().AddComponentToEntity(bgs.Entity(), engine.Default().Component(movecompname), mvc)
		}
	}

	s0 := engine.Default().NewSystem("", 0, func(ctx ecs.Context) {
		screen := ctx.World().Get("screen").(*ebiten.Image)
		fps := ebiten.CurrentFPS()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f fps", fps), 0, 0)
		ebitenutil.DebugPrintAt(screen, "d: toggle debug draw", 0, 15)
		ebitenutil.DebugPrintAt(screen, "p: toggle pause", 0, 30)
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			core.DebugDraw = !core.DebugDraw
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			movementPaused = !movementPaused
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyN) && movementPaused {
			xframes <- struct{}{}
		}
	})
	s0.AddTag(tau.WorldTagDraw)

}
