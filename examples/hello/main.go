// +build example

package main

import (
	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/examples/hello/hello"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  640,
		Height: 480,
		Scale:  1,
		Title:  "Hello, World!",
	})
	//
	w := engine.NewWorld(0)
	ecs.RegisterWorldDefaults(w)
	//
	e1 := w.NewEntity()
	hello.SetHelloComponentData(w, e1, hello.Hello{
		Text: "Hello",
		X:    100,
		Y:    100,
	})
	//
	e2 := w.NewEntity()
	hello.SetHelloComponentData(w, e2, hello.Hello{
		Text: "Primen",
		X:    120,
		Y:    120,
	})
	hello.SetMoveComponentData(w, e2, hello.Move{
		XSpeed: hello.SPEED,
		YSpeed: hello.SPEED,
	})
	_ = engine.Run()
}
