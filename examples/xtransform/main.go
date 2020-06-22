package main

import (
	"math"

	"github.com/gabstv/ecs/v2"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
)

func main() {
	core.DebugDraw = true
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     800,
		Height:    600,
		Resizable: true,
		OnReady:   ready,
	})
	engine.Run()
}

func ready(engine primen.Engine) {
	w := engine.NewWorldWithDefaults(0)
	tr := primen.NewRootNode(w)
	coretr := tr.Transform()
	coretr.SetX(400).SetY(300)
	tr2 := primen.NewChildNode(tr)
	coretr2 := tr2.Transform()
	coretr2.SetX(10).SetY(20)
	coretr.SetAngle(-math.Pi / 2)
	core.SetFunctionComponentData(w, tr.Entity(), core.Function{
		Update: func(ctx core.UpdateCtx, e ecs.Entity) {
			dd := core.GetTransformComponentData(w, e)
			dd.SetAngle(dd.Angle() + math.Pi*ctx.DT())
		},
	})
	//
	tr3 := primen.NewChildNode(tr2)
	coretr3 := tr3.Transform()
	coretr3.SetX(-40).SetY(50)
	//
	core.SetFunctionComponentData(w, tr2.Entity(), core.Function{
		Update: func(ctx core.UpdateCtx, e ecs.Entity) {
			dd := core.GetTransformComponentData(w, e)
			dd.SetAngle(dd.Angle() - .5*math.Pi*ctx.DT())
		},
	})
}
