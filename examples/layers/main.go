// +build example

package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/examples/layers/res"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	// "github.com/pkg/profile" // enable line 27
)

var movementPaused = false
var xframes = make(chan struct{}, 30)
var globalScale = 1.0
var radiusScale = 1.0
var wave1 = 1.0
var waver = 0.0

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
		OnReady: func(e *primen.Engine) {
			dogamesetup(ctx, e, spbgs, spfgs)
		},
	})
	engine.AddEventListener("act_of_nature", func(eventName string, e core.Event) {
		println("act of nature happened!")
		println(e.Data.(ecs.Entity))
		if globalScale != 1 {
			globalScale = 1
		} else {
			globalScale = -2
		}
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
	ChildSprite *primen.Sprite
	HueShift    bool
}

func dogamesetup(ctx context.Context, engine *primen.Engine, bgs, fgs []*ebiten.Image) {
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
			waver += dt
			wave1 = 1 + math.Cos(waver)
			if waver > 2*math.Pi {
				waver -= 2 * math.Pi
			}
			//
			for _, match := range ctx.System().View().Matches() {
				sprite := match.Components[spc].(*core.Sprite)
				transform := match.Components[trc].(*core.Transform)
				drawlayer := match.Components[dlc].(*core.DrawLayer)
				movecomp := match.Components[moc].(*OrbitalMovement)
				movecomp.R += movecomp.Speed * dt * globalScale * wave1
				xx := math.Cos(movecomp.R) * movecomp.Dx * radiusScale
				yy := math.Sin(movecomp.R) * movecomp.Dy * radiusScale
				transform.X = movecomp.Ox + xx
				transform.Y = movecomp.Oy + yy
				transform.Angle += dt * (math.Pi / 4) * movecomp.AngleR
				if rand.Float64() < 0.0001 {
					newlayer := rand.Intn(4)
					drawlayer.Layer = core.LayerIndex(newlayer)
					sprite.Image = bgs[newlayer]
					movecomp.ChildSprite.SetImage(fgs[newlayer])
					movecomp.ChildSprite.SetLayer(primen.Layer(newlayer))
					ctx.Engine().DispatchEvent("act_of_nature", match.Entity)
				}
				if movecomp.HueShift {
					sprite.SetColorHue(waver * 1.2)
				}
			}
		},
	}
	//
	_ = movecs.Components(engine.Default())
	core.SetupSystem(engine.Default(), movecs)

	rand.Seed(112358)

	root := primen.NewTransform(engine.Root(nil))
	root.SetX(320 / 2)
	root.SetY(240 / 2)

	root.UpsertFns(func(ctx core.Context, e ecs.Entity) {
		t := ctx.World().Component(core.CNTransform).Data(e).(*core.Transform)
		t.X = float64(ctx.Engine().Width() / 2)
		t.Y = float64(ctx.Engine().Height() / 2)
	}, nil, nil)

	for i := 0; i < 4; i++ {
		for j := 0; j < 20; j++ {
			//ri := rand.Intn(4)
			rl := rand.Intn(4)
			bgs := primen.NewSprite(root, bgs[rl], core.LayerIndex(rl))
			if rand.Float64() < 0.1 {
				bgs.SetColorHue(rand.Float64() * (math.Pi * 2))
			}
			if rand.Float64() < 0.06 {
				println("COMPOSITE LIGHTER")
				bgs.SetCompositeMode(ebiten.CompositeModeLighter)
			} else if rand.Float64() < 0.03 {
				println("COMPOSITE XOR")
				bgs.SetCompositeMode(ebiten.CompositeModeXor)
			}
			bgs.SetOrigin(.5, .5)
			fgs := primen.NewSprite(bgs, fgs[rl], core.LayerIndex(rl))
			fgs.SetOrigin(.5, .5)
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
			if rand.Float64() < 0.12 {
				println("HUE SHIFTER")
				mvc.HueShift = true
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
		if inpututil.IsKeyJustPressed(ebiten.KeyRightBracket) {
			radiusScale += 0.1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeftBracket) {
			radiusScale -= 0.1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			engine.SetScreenScale(ebiten.DeviceScaleFactor())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			engine.SetScreenScale(.5)
		}
	})
	s0.AddTag(primen.WorldTagDraw)

}
