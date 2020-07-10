package main

import (
	"io/ioutil"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core/ui/imgui"
	"github.com/gabstv/primen/dom"
)

func main() {
	fb, _ := ioutil.ReadFile("ui.xml")
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		OnReady: func(e primen.Engine) {
			node, _ := dom.ParseXMLText(string(fb))
			imgui.AddUI(node.(dom.ElementNode))
		},
	})
	imgui.Setup(engine)
	engine.Run()
}
