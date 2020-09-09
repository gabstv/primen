package main

import (
	"io/ioutil"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/modules/imgui"
)

func main() {
	fb, _ := ioutil.ReadFile("ui.xml")
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     800,
		Height:    600,
		Resizable: true,
		OnReady: func(e primen.Engine) {
			nodes, err := dom.ParseXMLString(string(fb))
			if err != nil {
				panic(err)
			}
			imgui.AddUI(nodes)
			e.AddEventListener("test", func(eventName string, e core.Event) {
				println(e.Data.(string))
			})
		},
	})
	engine.SetDebugTPS(true)
	imgui.Setup(engine)
	engine.Run()
}
