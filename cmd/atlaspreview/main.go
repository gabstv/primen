package main

import (
	"io/ioutil"
	"os"

	"github.com/gabstv/ecs"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/urfave/cli"
)

var (
	colorTitle = primen.ColorFromHex("#f1c40f")
	colorItem  = primen.ColorFromHex("ecf0f1")
	uiLayer    = primen.Layer2
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
			Scale:     ebiten.DeviceScaleFactor(),
		})
		return engine.Run()
	}
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func buildReady(c *cli.Context) func(e *primen.Engine) {
	core.DebugDraw = true
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
		_ = loadAtlas(e, ff)
	}
}

func loadAtlas(e *primen.Engine, atlas *io.Atlas) *AtlasPreviewer {
	//TODO: remove older one if present (to allow load multiple)
	return newAtlasPreviewer(e, atlas)
}

func errready(v string) func(e *primen.Engine) {
	return func(e *primen.Engine) {
		primen.SetDrawFuncs(e.Default(), e.Default().NewEntity(), nil, func(ctx core.Context, e ecs.Entity) {
			ebitenutil.DebugPrint(ctx.Screen(), v)
		}, nil)
	}
}

type AtlasPreviewer struct {
	e        *primen.Engine
	itemList *AtlasItemList
	atlas    *io.Atlas
	canvas   *primen.Transform
	titler   *primen.Transform
}

func (p *AtlasPreviewer) Destroy() {
	p.itemList.Destroy()
	p.itemList = nil
}

func (p *AtlasPreviewer) resetCanvas() {
	if p.canvas != nil {
		primen.Destroy(p.canvas)
	}
	p.canvas = nil
	if p.titler != nil {
		primen.Destroy(p.titler)
	}
	p.titler = nil
}

func (p *AtlasPreviewer) createCanvas() {
	p.canvas = primen.NewTransform(p.e.Root(nil))
	p.canvas.SetPos(float64(p.e.Width())/2, float64(p.e.Height())/2)
	p.titler = primen.NewTransform(p.e.Root(nil))
	p.titler.SetPos(float64(p.e.Width())/2, 30)
}

func (p *AtlasPreviewer) setupAnim(name string) {
	p.resetCanvas()
	p.createCanvas()
	p.canvas.SetScale2(2 * ebiten.DeviceScaleFactor())
	anim := p.atlas.GetAnimation(name)
	title := primen.NewLabel(p.titler, nil, uiLayer)
	title.SetText("Animation: " + name)
	title.SetColor(primen.ColorFromHex("#e74c3c"))
	title.SetOrigin(.5, .5)
	as := primen.NewAnimatedSprite(p.canvas, primen.Layer0, anim)
	as.PlayClipIndex(0)
	// nclips := anim.Count()
	// anim.Each(func(i int, clip core.AnimationClip) bool {
	//
	// })
}

func newAtlasPreviewer(e *primen.Engine, atlas *io.Atlas) *AtlasPreviewer {
	p := &AtlasPreviewer{
		e:     e,
		atlas: atlas,
	}
	p.itemList = newAtlasItemList(p)
	if len(atlas.GetAnimations()) > 0 {
		p.setupAnim(atlas.GetAnimations()[0].Name)
	}
	return p
}

type AtlasItemList struct {
	tr *primen.Transform
}

func (al *AtlasItemList) Destroy() {
	primen.Destroy(al.tr)
	al.tr = nil
}

func newAtlasItemList(parent *AtlasPreviewer) *AtlasItemList {
	atlas := parent.atlas
	pp := primen.NewTransform(parent.e.Root(nil))
	pp.SetPos(10, 10)
	nexty := 0.0
	// animations
	{
		lbl := primen.NewLabel(pp, nil, uiLayer)
		lbl.SetColor(colorTitle)
		lbl.SetText("Animations:")
		lbl.SetY(nexty)
		_, yp := lbl.ComputedSize()
		nexty += float64(yp) + 10
		for _, animg := range atlas.GetAnimations() {
			li := primen.NewLabel(pp, nil, uiLayer)
			li.SetColor(colorItem)
			li.SetY(nexty)
			li.SetText(animg.Name)
			_, yp = li.ComputedSize()
			nexty += float64(yp) + 10
		}
	}
	return &AtlasItemList{
		tr: pp,
	}
}
