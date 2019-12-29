package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	// app.Commands = []cli.Command{
	// 	cli.Command{
	// 		Name:      "new",
	// 		ShortName: "n",
	// 		Action:    cmdNew,
	// 		Flags: []cli.Flag{
	// 			cli.StringFlag{
	// 				Name:   "package, p",
	// 				EnvVar: "GOPACKAGE",
	// 			},
	// 			cli.StringFlag{
	// 				Name: "component, c",
	// 			},
	// 			cli.IntFlag{
	// 				Name: "priority",
	// 			},
	// 		},
	// 	},
	// }
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "input, i",
		},
		cli.Float64Flag{
			Name:  "scale, s",
			Value: 2,
		},
	}

	app.Action = run

	if err := app.Run(os.Args); err != nil {
		if e, ok := err.(*cli.ExitError); ok {
			os.Exit(e.ExitCode())
		}
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	inputfn := c.String("input")
	eimg, img, err := ebitenutil.NewImageFromFile(inputfn, ebiten.FilterNearest)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	scale := c.Float64("scale")
	return ebiten.Run(ebrun(eimg, scale), w, h, scale, "XYZ")
}

func ebrun(baseimg *ebiten.Image, scale float64) func(screen *ebiten.Image) error {
	imo := ebiten.DrawImageOptions{}
	xc := color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 100,
	}
	yc := color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 100,
	}
	return func(screen *ebiten.Image) error {
		mx, my := ebiten.CursorPosition()
		maxx, maxy := screen.Bounds().Dx(), screen.Bounds().Dy()
		//
		if ebiten.IsDrawingSkipped() {
			return nil
		}
		screen.Fill(color.Black)
		//
		screen.DrawImage(baseimg, &imo)
		//
		ebitenutil.DrawLine(screen, 0, float64(my), float64(maxx), float64(my), xc)
		ebitenutil.DrawLine(screen, float64(mx), 0, float64(mx), float64(maxy), yc)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("mx: %v my: %v", mx, my))
		return nil
	}
}
