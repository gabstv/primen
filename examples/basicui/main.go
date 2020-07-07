package main

import (
	"io/ioutil"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:   800,
		Height:  600,
		Scale:   1,
		OnReady: ready,
	})
	engine.Run()
}

func ready(e primen.Engine) {
	fb, _ := ioutil.ReadFile("ui.xml")
	root, err := dom.ParseXMLText(string(fb))
	if err != nil {
		panic(err)
	}
	w := e.NewWorldWithDefaults(0)
	entity := w.NewEntity()
	core.SetUIManagerComponentData(w, entity, core.NewUIManager())
	core.GetUIManagerComponentData(w, entity).Setup(root.(dom.ElementNode))
}
