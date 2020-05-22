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
	"path/filepath"

	"github.com/gabstv/primen/internal/spriteutils"
	"github.com/gabstv/primen/spr"
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
	_, sfn := filepath.Split(inputfn)
	//
	return ebiten.Run(ebrun(ebruninput{
		Source: eimg,
		Scale:  scale,
		Def:    def,
		Meta:   spriteutils.MayGetSpriteDef(def.Metadata),
	}), w, h, scale, "SPR CENTER - "+sfn)
}

var xc = color.RGBA{
	R: 255,
	G: 0,
	B: 0,
	A: 100,
}
var yc = color.RGBA{
	R: 0,
	G: 255,
	B: 0,
	A: 100,
}
var prevc = color.RGBA{
	R: 255,
	G: 200,
	B: 40,
	A: 100,
}

type runctx struct {
	In           ebruninput
	Opt          *ebiten.DrawImageOptions
	Ox           int64
	Oy           int64
	Oxf          float64
	Oyf          float64
	Xoff         float64
	Yoff         float64
	Offstate     byte
	Silent       bool
	HelpOn       bool
	CursorX      int
	CursorY      int
	ScreenWidth  int
	ScreenHeight int
	ImageWidth   int
	ImageHeight  int
}

func ebupdate(screen *ebiten.Image, c *runctx) {
	ebinput(screen, c)
}

func ebinput(screen *ebiten.Image, c *runctx) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		c.In.Meta.ExtraOriginX = c.Xoff
		c.In.Meta.ExtraOriginY = c.Yoff
		c.In.Meta.RawOriginX = int64(c.CursorX)
		c.In.Meta.RawOriginY = int64(c.CursorY)
		c.Ox = int64(c.CursorX)
		c.Oy = int64(c.CursorY)
		c.Oxf = float64(c.Ox) + c.In.Meta.ExtraOriginX
		c.Oyf = float64(c.Oy) + c.In.Meta.ExtraOriginY
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
		c.Offstate++
		if c.Offstate > 2 {
			c.Offstate = 0
		}
		switch c.Offstate {
		case 0:
			c.Xoff = 0.5
			c.Yoff = 0.5
		case 1:
			c.Xoff = 1
			c.Yoff = 1
		case 2:
			c.Xoff = 0
			c.Yoff = 0
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyH) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			c.Silent = !c.Silent
		} else {
			c.HelpOn = !c.HelpOn
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		d := c.In.Def
		d.Metadata = json.RawMessage(c.In.Meta.MustJSON())
		d.Size = spr.Vec2{
			X: float64(c.ImageWidth),
			Y: float64(c.ImageHeight),
		}
		d.Origin = spr.Vec2{
			X: c.Oxf,
			Y: c.Oyf,
		}
		d.WriteToFile(d.Filename(), 0744)
		os.Exit(0)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) || inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		os.Exit(0)
	}
}

func ebdraw(screen *ebiten.Image, c *runctx) {
	screen.Fill(color.Black)
	//
	screen.DrawImage(c.In.Source, c.Opt)
	// current origin:
	ebitenutil.DrawLine(screen, 0, float64(c.Oy), float64(c.ScreenWidth), float64(c.Oy), prevc)
	ebitenutil.DrawLine(screen, float64(c.Ox), 0, float64(c.Ox), float64(c.ScreenHeight), prevc)
	// mouse position
	ebitenutil.DrawLine(screen, 0, float64(c.CursorY), float64(c.ScreenWidth), float64(c.CursorY), xc)
	ebitenutil.DrawLine(screen, float64(c.CursorX), 0, float64(c.CursorX), float64(c.ScreenHeight), yc)
	// ~ ~ ~ text
	// mouse pos
	if !c.Silent {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Size: (%v; %v); [h]elp", c.ImageWidth, c.ImageHeight), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: (%v; %v)", c.CursorX, c.CursorY), 0, 14)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: (%v; %v)", c.Xoff, c.Yoff), 0, 14*2)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Origin: (%v; %v)", c.Oxf, c.Oyf), 0, 14*3)
		if c.HelpOn {
			ebitenutil.DebugPrintAt(screen, "[H]elp:", 0, 14*4)
			ebitenutil.DebugPrintAt(screen, "[lmouse]: set origin;", 0, 14*5)
			ebitenutil.DebugPrintAt(screen, "[tab]: switch px offset;", 0, 14*6)
			ebitenutil.DebugPrintAt(screen, "[return]: save and quit", 0, 14*7)
			ebitenutil.DebugPrintAt(screen, "[q; esc]: quit w/o saving", 0, 14*8)
		}
	} else {
		if c.HelpOn {
			ebitenutil.DebugPrintAt(screen, "[H]elp:", 0, 14*0)
			ebitenutil.DebugPrintAt(screen, "[lmouse]: set origin;", 0, 14*1)
			ebitenutil.DebugPrintAt(screen, "[tab]: switch px offset;", 0, 14*2)
			ebitenutil.DebugPrintAt(screen, "[return]: save and quit", 0, 14*3)
			ebitenutil.DebugPrintAt(screen, "[q; esc]: quit w/o saving", 0, 14*4)
		}
	}
}

func ebrun(input ebruninput) func(screen *ebiten.Image) error {
	rctx := &runctx{
		In:          input,
		Opt:         &ebiten.DrawImageOptions{},
		Ox:          input.Meta.RawOriginX,
		Oy:          input.Meta.RawOriginY,
		Oxf:         float64(input.Meta.RawOriginX) + input.Meta.ExtraOriginX,
		Oyf:         float64(input.Meta.RawOriginY) + input.Meta.ExtraOriginY,
		Xoff:        0.5,
		Yoff:        0.5,
		ImageWidth:  input.Source.Bounds().Dx(),
		ImageHeight: input.Source.Bounds().Dy(),
	}
	return func(screen *ebiten.Image) error {
		mx, my := ebiten.CursorPosition()
		maxx, maxy := screen.Bounds().Dx(), screen.Bounds().Dy()
		rctx.CursorX = mx
		rctx.CursorY = my
		rctx.ScreenWidth = maxx
		rctx.ScreenHeight = maxy
		//
		ebupdate(screen, rctx)
		//
		if ebiten.IsDrawingSkipped() {
			return nil
		}
		ebdraw(screen, rctx)
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
