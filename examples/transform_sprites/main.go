package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var randomsprites [][]int

func init() {
	rand.Seed(time.Now().Unix())
	randomsprites = [][]int{
		[]int{0, 0, 16, 16},
		[]int{0, 16, 16, 32},
		[]int{16, 0, 32, 16},
		[]int{16, 16, 32, 32},
	}
}

type spinner struct {
	Speed float64
}

var spinnercs = &core.BasicCS{
	SysName: "spinnercs",
	SysExec: func(ctx core.Context) {
		sys := ctx.System()
		dt := ctx.DT()
		view := sys.View()
		tc := ctx.World().Component(core.CNTransform)
		sc := ctx.World().Component("spinnercs")
		scaleadd := sys.Get("scaleadd").(float64) + dt
		sys.Set("scaleadd", scaleadd)
		//
		xs := float64(0)
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			xs = 50
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			xs = -50
		}
		ys := float64(0)
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			ys = -50
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			ys = 50
		}
		rs := float64(0)
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			rs = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			rs = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			rs *= 0.25
		}
		for _, v := range view.Matches() {
			spin := v.Components[sc].(*spinner)
			tr := v.Components[tc].(*core.Transform)
			tr.Angle += spin.Speed * dt * rs
			tr.X += xs * dt
			tr.Y += ys * dt
			tr.ScaleY = 0.7 + math.Cos(scaleadd)/4
		}
	},
	GetComponents: func(w *ecs.World) []*ecs.Component {
		return []*ecs.Component{
			w.Component(core.CNTransform),
			core.UpsertComponent(w, ecs.NewComponentInput{
				Name: "spinnercs",
			}),
		}
	},
	SysInit: func(w *ecs.World, sys *ecs.System) {
		sys.Set("scaleadd", float64(0))
	},
}

func main() {
	ebimg, _, _ := ebitenutil.NewImageFromFile("img.png", ebiten.FilterDefault)

	engine := primen.NewEngine(&primen.NewEngineInput{
		Title:  "Basic Transform With Sprites",
		Width:  640,
		Height: 480,
		Scale:  2,
	})

	core.DebugDraw = true

	dw := engine.Default()
	sc := dw.Component(core.CNDrawable)
	tc := dw.Component(core.CNTransform)
	spinnercomp := spinnercs.Components(dw)[1]
	e := dw.NewEntity()
	t99 := &core.Transform{
		X:      320 / 2,
		Y:      240 / 2,
		ScaleX: 0.5,
		ScaleY: 0.7,
	}
	dw.AddComponentToEntity(e, tc, t99)
	dw.AddComponentToEntity(e, spinnercomp, &spinner{
		Speed: 3,
	})
	// add children
	for i := 0; i < 10; i++ {
		e2 := dw.NewEntity()
		//mm := primen.IM.Moved(primen.V(30, 0)).Rotated(primen.ZV, (math.Pi*2)*(float64(i)/10)).Project(primen.ZV)

		dw.AddComponentToEntity(e2, tc, &core.Transform{
			X:      math.Cos((math.Pi*2)*(float64(i)/10)) * 30,
			Y:      math.Sin((math.Pi*2)*(float64(i)/10)) * 30,
			Parent: t99,
			ScaleX: 1,
			ScaleY: 1,
		})
		ri := randomsprites[rand.Intn(4)]
		dw.AddComponentToEntity(e2, sc, &core.Sprite{
			Bounds:  image.Rect(ri[0], ri[1], ri[2], ri[3]),
			Image:   ebimg,
			ScaleX:  1,
			ScaleY:  1,
			OriginX: 0.5,
			OriginY: 0.5,
		})
	}
	core.SetupSystem(dw, spinnercs)

	// debug system
	ddrawsys := dw.NewSystem("", -100, func(ctx ecs.Context) {
		screen := ctx.World().Get("screen").(*ebiten.Image)
		fps := ebiten.CurrentFPS()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f fps", fps), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("x = %.2f; y = %.2f;", t99.X, t99.Y), 0, 12)
		ebitenutil.DebugPrintAt(screen, "arrows = move; x, z = rotate", 0, 24)

	}, spinnercomp)
	ddrawsys.AddTag(primen.WorldTagDraw)

	engine.Run()
}
