package main

import (
	"time"

	"github.com/gabstv/primen"
	osfs "github.com/gabstv/primen/io/os"
)

func main() {
	qch := make(chan struct{})
	e := primen.NewEngine(&primen.NewEngineInput{
		Width:  400,
		Height: 300,
		FS:     osfs.New("../shared"),
		OnReady: func(e primen.Engine) {
			go func() {
				time.Sleep(time.Millisecond * 5500)
				e.Exit()
			}()
			go func() {
				<-e.Ctx().Done()
				close(qch)
			}()
		},
	})
	e.Run()
	<-qch
}
