package main

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var engine *groove.Engine
var hellocomp *ecs.Component

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
		Name: "init_engine",
	})
	if err != nil {
		panic(err)
	}
	hellocomp = comp
	sys0 := world.NewSystem(0, initEngineSystemExec, comp)
	sys0.AddTag(groove.WorldTagDraw)
	entity0 := world.NewEntity()
	world.AddComponentToEntity(entity0, comp, &initEngineData{"Hello,", 30, 40})
	entity1 := world.NewEntity()
	world.AddComponentToEntity(entity1, comp, &initEngineData{"World!", 50, 60})
	// run
	engine.Run()
}

type initEngineData struct {
	Text string
	X    int
	Y    int
}

func initEngineSystemExec(dt float64, view *ecs.View) {
	img := engine.Get(groove.EbitenScreen).(*ebiten.Image)
	for _, v := range view.Matches() {
		data := v.Components[hellocomp].(*initEngineData)
		ebitenutil.DebugPrintAt(img, data.Text, data.X, data.Y)
	}
}
