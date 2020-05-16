package main

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var engine *tau.Engine

func hellocomp(w *ecs.World) *ecs.Component {
	return tau.UpsertComponent(w, ecs.NewComponentInput{
		Name: "hellocs_comp",
	})
}

func movecomp(w *ecs.World) *ecs.Component {
	return tau.UpsertComponent(w, ecs.NewComponentInput{
		Name: "movecs_comp",
	})
}

var hellocs = &tau.BasicCS{
	SysName: "hellocs_system",
	SysExec: initEngineSystemExec,
	SysTags: []string{tau.WorldTagDraw},
	GetComponents: func(w *ecs.World) []*ecs.Component {
		return []*ecs.Component{
			hellocomp(w),
		}
	},
}

var movecs = &tau.BasicCS{
	SysName: "movecs_system",
	SysExec: moveSysExec,
	SysTags: []string{tau.WorldTagUpdate},
	GetComponents: func(w *ecs.World) []*ecs.Component {
		return []*ecs.Component{
			movecomp(w),
			hellocomp(w),
		}
	},
}

const SPEED float64 = 120

func main() {
	engine = tau.NewEngine(&tau.NewEngineInput{
		Width:  320,
		Height: 240,
		Scale:  2,
		Title:  "Hello, World!",
	})
	// add components and systems
	world := engine.Default()
	tau.SetupSystem(world, hellocs)
	tau.SetupSystem(world, movecs)

	entity0 := world.NewEntity()
	world.AddComponentToEntity(entity0, hellocomp(world), &initEngineData{"Hello,", 30, 40})
	entity1 := world.NewEntity()
	world.AddComponentToEntity(entity1, hellocomp(world), &initEngineData{"World!", 50, 60})
	world.AddComponentToEntity(entity1, movecomp(world), &moveCompData{
		XSpeed: SPEED,
		YSpeed: SPEED,
	})
	// run
	engine.Run()
}

type initEngineData struct {
	Text string
	X    int
	Y    int
}

type moveCompData struct {
	XSpeed float64
	YSpeed float64
	XSum   float64
	YSum   float64
}

func initEngineSystemExec(ctx tau.Context) {
	screen := ctx.Screen()
	c := hellocomp(ctx.World())
	for _, v := range ctx.System().View().Matches() {
		data := v.Components[c].(*initEngineData)
		ebitenutil.DebugPrintAt(screen, data.Text, data.X, data.Y)
	}
}

func moveSysExec(ctx tau.Context) {
	dt := ctx.DT()
	helloc := hellocomp(ctx.World())
	movec := movecomp(ctx.World())
	for _, v := range ctx.System().View().Matches() {
		iedata := v.Components[helloc].(*initEngineData)
		movedata := v.Components[movec].(*moveCompData)
		movedata.XSum += dt * movedata.XSpeed
		movedata.YSum += dt * movedata.YSpeed
		for movedata.XSum >= 1 {
			iedata.X++
			movedata.XSum--
		}
		for movedata.XSum <= -1 {
			iedata.X--
			movedata.XSum++
		}
		for movedata.YSum >= 1 {
			iedata.Y++
			movedata.YSum--
		}
		for movedata.YSum <= -1 {
			iedata.Y--
			movedata.YSum++
		}
		if iedata.X >= 280 && movedata.XSpeed > 0 {
			movedata.XSpeed = -SPEED
		}
		if iedata.X <= 0 && movedata.XSpeed < 0 {
			movedata.XSpeed = SPEED
		}
		if iedata.Y >= 220 && movedata.YSpeed > 0 {
			movedata.YSpeed = -SPEED
		}
		if iedata.Y <= 0 && movedata.YSpeed < 0 {
			movedata.YSpeed = SPEED
		}
	}
}
