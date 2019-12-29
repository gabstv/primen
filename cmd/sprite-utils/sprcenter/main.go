package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/gabstv/troupe/internal/spriteutils"
	"github.com/gabstv/troupe/spr"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
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

type ebruninput struct {
	Source *ebiten.Image
	Scale  float64
	Def    *spr.SpriteDef
	Meta   *spriteutils.SpriteDefMeta
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
	//
	def, _, err := loadOrCreateDef(inputfn)
	if err != nil {
		return cli.NewExitError(err.Error(), 3)
	}
	//
	return ebiten.Run(ebrun(ebruninput{
		Source: eimg,
		Scale:  scale,
		Def:    def,
		Meta:   spriteutils.MayGetSpriteDef(def.Metadata),
	}), w, h, scale, "XYZ")
}

func ebrun(input ebruninput) func(screen *ebiten.Image) error {
	baseimg := input.Source
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
	prevc := color.RGBA{
		R: 255,
		G: 200,
		B: 40,
		A: 100,
	}
	ox := input.Meta.RawOriginX
	oy := input.Meta.RawOriginY
	oxf := float64(ox) + input.Meta.ExtraOriginX
	oyf := float64(oy) + input.Meta.ExtraOriginY
	return func(screen *ebiten.Image) error {
		mx, my := ebiten.CursorPosition()
		maxx, maxy := screen.Bounds().Dx(), screen.Bounds().Dy()
		//
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			input.Meta.ExtraOriginX = 0.5
			input.Meta.ExtraOriginY = 0.5
			input.Meta.RawOriginX = int64(mx)
			input.Meta.RawOriginY = int64(my)
			ox = int64(mx)
			oy = int64(my)
			oxf = float64(ox) + input.Meta.ExtraOriginX
			oyf = float64(oy) + input.Meta.ExtraOriginY
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
			d := input.Def
			d.Metadata = json.RawMessage(input.Meta.MustJSON())
			d.Size = spr.Vec2{
				X: float64(maxx),
				Y: float64(maxy),
			}
			d.Origin = spr.Vec2{
				X: oxf,
				Y: oyf,
			}
			d.WriteToFile(d.Filename(), 0744)
			os.Exit(0)
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyEscape) || inpututil.IsKeyJustReleased(ebiten.KeyQ) {
			os.Exit(0)
		}
		//
		if ebiten.IsDrawingSkipped() {
			return nil
		}
		screen.Fill(color.Black)
		//
		screen.DrawImage(baseimg, &imo)
		// current origin:
		ebitenutil.DrawLine(screen, 0, float64(oy), float64(maxx), float64(oy), prevc)
		ebitenutil.DrawLine(screen, float64(ox), 0, float64(ox), float64(maxy), prevc)
		// mouse position
		ebitenutil.DrawLine(screen, 0, float64(my), float64(maxx), float64(my), xc)
		ebitenutil.DrawLine(screen, float64(mx), 0, float64(mx), float64(maxy), yc)
		// ~ ~ ~ text
		// mouse pos
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Size: (%v; %v)", maxx, maxy), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: (%v; %v)", mx, my), 0, 14)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Origin: (%v; %v)", oxf, oyf), 0, 14*2)
		return nil
	}
}

var zerometa = &spriteutils.SpriteDefMeta{
	RawOriginX:   0,
	RawOriginY:   0,
	ExtraOriginX: 0.5,
	ExtraOriginY: 0.5,
}

func loadOrCreateDef(imgname string) (*spr.SpriteDef, bool, error) {
	if _, err := os.Stat(imgname + ".spr"); err == nil {
		if sd, err := spr.ReadSpriteDefFile(imgname + ".spr"); err == nil {
			return sd, false, nil
		}
	}
	if _, err := os.Stat(imgname + ".spr.meta"); err == nil {
		if sd, err := spr.ReadSpriteDefFile(imgname + ".spr.meta"); err == nil {
			return sd, false, nil
		}
	}
	if _, err := os.Stat(imgname + ".meta"); err == nil {
		if sd, err := spr.ReadSpriteDefFile(imgname + ".meta"); err == nil {
			return sd, false, nil
		}
	}
	nd, err := newDef(imgname)
	if err != nil {
		return nil, false, err
	}
	return nd, true, nil
}

func newDef(imgname string) (*spr.SpriteDef, error) {
	f, err := os.Open(imgname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	def := spr.NewSpriteDef(imgname+".meta", zerometa.MustJSON())
	def.Source = &spr.SourceData{
		File:   imgname,
		Width:  float64(i.Bounds().Dx()),
		Height: float64(i.Bounds().Dy()),
	}
	def.Size = spr.Vec2{
		X: float64(i.Bounds().Dx()),
		Y: float64(i.Bounds().Dy()),
	}
	return def, nil
}
