package main

import (
	"io/ioutil"
	"os"

	"github.com/gabstv/ecs/v2"
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

func buildReady(c *cli.Context) func(e primen.Engine) {
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
	return func(e primen.Engine) {
		_ = loadAtlas(e, ff)
	}
}

func loadAtlas(e primen.Engine, atlas *io.Atlas) *AtlasPreviewer {
	//TODO: remove older one if present (to allow load multiple)
	return newAtlasPreviewer(e, atlas)
}

func errready(v string) func(e primen.Engine) {
	return func(e primen.Engine) {
		w := e.NewWorldWithDefaults(0)
		f := primen.NewRootFnNode(w)
		f.Function().Draw = func(ctx core.DrawCtx, e ecs.Entity) {
			ebitenutil.DebugPrint(ctx.Renderer().Screen(), v)
		}
	}
}

type AtlasPreviewer struct {
	e        primen.Engine
	itemList *AtlasItemList
	atlas    *io.Atlas
	canvas   *primen.Node
	titler   *primen.Node
	w        primen.World
}

func (p *AtlasPreviewer) Destroy() {
	p.itemList.Destroy()
	p.itemList = nil
}

func (p *AtlasPreviewer) resetCanvas() {
	if p.canvas != nil {
		p.canvas.Destroy()
	}
	p.canvas = nil
	if p.titler != nil {
		p.titler.Destroy()
	}
	p.titler = nil
}

func (p *AtlasPreviewer) createCanvas() {
	p.canvas = primen.NewRootNode(p.w)
	p.canvas.Transform().SetX(float64(p.e.Width()) / 2).SetY(float64(p.e.Height()) / 2)
	p.titler = primen.NewRootNode(p.w)
	p.titler.Transform().SetX(float64(p.e.Width()) / 2).SetY(30)
}

func (p *AtlasPreviewer) setupAnim(name string) {
	p.resetCanvas()
	p.createCanvas()
	p.canvas.Transform().SetScale(8*ebiten.DeviceScaleFactor(), 8*ebiten.DeviceScaleFactor())
	anim := p.atlas.GetAnimation(name)
	title := primen.NewChildLabelNode(p.titler, uiLayer)
	title.Label().SetText("Animation: " + name).SetColor(primen.ColorFromHex("#e74c3c"))
	title.Label().SetOrigin(.5, .5)
	as := primen.NewChildAnimatedSpriteNode(p.canvas, primen.Layer0, 12, anim)
	as.SpriteAnim().PlayClipIndex(0)
	// nclips := anim.Count()
	// anim.Each(func(i int, clip core.AnimationClip) bool {
	//
	// })
}

func newAtlasPreviewer(e primen.Engine, atlas *io.Atlas) *AtlasPreviewer {
	p := &AtlasPreviewer{
		e:     e,
		atlas: atlas,
		w:     e.NewWorldWithDefaults(0),
	}
	p.itemList = newAtlasItemList(p)
	if len(atlas.GetAnimations()) > 0 {
		p.setupAnim(atlas.GetAnimations()[0].Name)
	}
	return p
}

type AtlasItemList struct {
	tr *primen.Node
}

func (al *AtlasItemList) Destroy() {
	al.tr.Destroy()
	al.tr = nil
}

func newAtlasItemList(parent *AtlasPreviewer) *AtlasItemList {
	atlas := parent.atlas
	//pp := primen.NewTransform(parent.e.Root(nil))
	pp := primen.NewRootNode(parent.w)
	pp.Transform().SetX(10).SetY(10)
	nexty := 0.0
	// animations
	{
		//lbl := primen.NewLabel(pp, nil, uiLayer)
		lbl := primen.NewChildLabelNode(pp, uiLayer)
		lbl.Label().SetColor(colorTitle)
		lbl.Label().SetText("Animations:")
		lbl.Transform().SetY(nexty)
		p := lbl.Label().ComputedSize()
		nexty += float64(p.Y) + 10
		for _, animg := range atlas.GetAnimations() {
			li := primen.NewChildLabelNode(pp, uiLayer)
			li.Label().SetColor(colorItem)
			li.Transform().SetY(nexty)
			li.Label().SetText(animg.Name)
			p = li.Label().ComputedSize()
			nexty += float64(p.Y) + 10
		}
	}
	return &AtlasItemList{
		tr: pp,
	}
}
