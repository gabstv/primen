package main

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var engine *groove.Engine
var hellocomp *ecs.Component
var movecomp *ecs.Component

const SPEED float64 = 120

func main() {
	engine = groove.NewEngine(&groove.NewEngineInput{
		Width:  320,
		Height: 240,
		Scale:  2,
		Title:  "Hello, World!",
	})
	// add components and systems
	world := engine.Default()
	comp, err := world.NewComponent(ecs.NewComponentInput{
		Name: "hello",
	})
	if err != nil {
		panic(err)
	}
	hellocomp = comp
	movecomp, err = world.NewComponent(ecs.NewComponentInput{
		Name: "move",
	})
	if err != nil {
		panic(err)
	}
	sys0 := world.NewSystem(0, initEngineSystemExec, hellocomp)
	sys0.AddTag(groove.WorldTagDraw)
	sys1 := world.NewSystem(1, moveSysExec, movecomp, hellocomp)
	sys1.AddTag(groove.WorldTagUpdate)
	entity0 := world.NewEntity()
	world.AddComponentToEntity(entity0, hellocomp, &initEngineData{"Hello,", 30, 40})
	entity1 := world.NewEntity()
	world.AddComponentToEntity(entity1, hellocomp, &initEngineData{"World!", 50, 60})
	world.AddComponentToEntity(entity1, movecomp, &moveCompData{
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

func initEngineSystemExec(dt float64, view *ecs.View, sys *ecs.System) {
	img := engine.Get(groove.EbitenScreen).(*ebiten.Image)
	for _, v := range view.Matches() {
		data := v.Components[hellocomp].(*initEngineData)
		ebitenutil.DebugPrintAt(img, data.Text, data.X, data.Y)
	}
}

func moveSysExec(dt float64, view *ecs.View, sys *ecs.System) {
	for _, v := range view.Matches() {
		iedata := v.Components[hellocomp].(*initEngineData)
		movedata := v.Components[movecomp].(*moveCompData)
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
