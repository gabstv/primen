package main

import (
	"io/ioutil"
	"os"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "atlaspreview"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "width",
			Value: 800,
		},
		cli.IntFlag{
			Name:  "height",
			Value: 600,
		},
	}
	app.Action = func(c *cli.Context) error {
		engine := primen.NewEngine(&primen.NewEngineInput{
			Width:     c.Int("width"),
			Height:    c.Int("height"),
			Resizable: true,
			OnReady:   buildReady(c),
			Title:     "PRIMEN - Atlas Preview",
		})
		return engine.Run()
	}
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func buildReady(c *cli.Context) func(e *primen.Engine) {
	fn := c.Args().First()
	if fn == "" {
		return errready("No atlas file specified")
	}
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return errready(err.Error())
	}
	ff, err := io.ParseAtlas(b)
	if err != nil {
		return errready(err.Error())
	}
	println(ff)
	return func(e *primen.Engine) {
		println("hey")
		for _, a := range ff.GetAnimations() {
			println(a.Name)
		}
	}
}

func errready(v string) func(e *primen.Engine) {
	return func(e *primen.Engine) {
		primen.SetDrawFuncs(e.Default(), e.Default().NewEntity(), nil, func(ctx core.Context, e ecs.Entity) {
			ebitenutil.DebugPrint(ctx.Screen(), v)
		}, nil)
	}
}
