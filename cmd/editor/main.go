package main

import (
	"os"

	"github.com/gabstv/primen"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "width",
			Value:  800,
			EnvVar: "TAU_EDITOR_WIDTH",
			Usage:  "Initial width",
		},
		cli.IntFlag{
			Name:   "height",
			Value:  600,
			EnvVar: "TAU_EDITOR_HEIGHT",
			Usage:  "Initial height",
		},
	}
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	e := primen.NewEngine(&primen.NewEngineInput{
		Width:     c.Int("width"),
		Height:    c.Int("height"),
		Scale:     1,
		Title:     "Tau Editor",
		Resizable: true,
	})
	return e.Run()
}
