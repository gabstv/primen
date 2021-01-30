package main

import (
	"math"
	"os"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	paudio "github.com/gabstv/primen/audio"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func main() {
	f, _ := os.Open("saw.wav")
	ws, err := wav.Decode(paudio.Context(), f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	panwave := paudio.NewStereoPanStreamFromReader(ws)
	panloop := audio.NewInfiniteLoop(panwave, ws.Length())
	plr, _ := audio.NewPlayer(paudio.Context(), panloop)
	plr.SetVolume(.5)
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		Title:  "PAN TEST",
		OnReady: func(e primen.Engine) {
			plr.Play()
			w := e.NewWorldWithDefaults(0)
			rfn := primen.NewRootFnNode(w)
			rw := 0.0
			rfn.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
				rw += ctx.DT()
				panwave.SetPan(math.Cos(rw))
			}
		},
	})
	engine.Run()
}
