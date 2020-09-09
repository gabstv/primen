package main

import (
	"image/color"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/modules/imgui"
	"github.com/hajimehoshi/ebiten"
	im "github.com/inkyblackness/imgui-go/v2"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     1024,
		Height:    768,
		Resizable: true,
		Scale:     1,
		OnReady:   ready,
		Title:     "Custom UI",
	})
	// setup imgui module
	imgui.Setup(engine)
	engine.Run()
}

func ready(e primen.Engine) {
	// https://github.com/ocornut/imgui/blob/53f0f972737a42472fce671864cb9b5fa4514562/docs/FONTS.md
	fcfg := im.NewFontConfig()
	fcfg.SetSize(18)
	//fcfg.SetMergeMode(false)
	fcfg.SetOversampleV(3)
	fcfg.SetOversampleH(3)
	fcfg.SetPixelSnapH(true)
	//fcfg.SetGlyphMaxAdvanceX(18)
	//fcfg.SetGlyphMinAdvanceX(18)
	gr := im.CurrentIO().Fonts().GlyphRangesDefault()
	if err := im.CurrentIO().Fonts().BuildWithFreeType(); err != nil {
		panic(err)
	}
	uidata.f = im.CurrentIO().Fonts().AddFontFromFileTTFV("./nasalization-rg.ttf", 18, fcfg, gr)
	imz := im.CurrentIO().Fonts().TextureDataRGBA32() // this needs to be called so the new loaded font(s) is/are packed
	uidata.xx, uidata.yy = imz.Width, imz.Height
	im.CurrentIO().SetFontGlobalScale(1)
	imgui.SetFilter(ebiten.FilterNearest)
	uidata.id = imgui.AddRawUI(customUI)

	// for true pixel perfect fonts:
	// fcfg.SetSize(13)
	// //fcfg.SetMergeMode(false)
	// fcfg.SetOversampleV(1)
	// fcfg.SetOversampleH(1)
	// fcfg.SetPixelSnapH(true)
}

var uidata = struct {
	open bool
	id   imgui.UID
	f    im.Font
	xx   int
	yy   int
}{
	open: true,
}

func customUI(ctx core.DrawCtx) {
	im.BeginV("Hello UI", &uidata.open, im.WindowFlagsNoCollapse)
	im.PushFont(uidata.f)
	im.PushStyleColor(im.StyleColorButton, imgui.Color(color.RGBA{
		R: 200,
		G: 50,
		B: 50,
		A: 200,
	}))
	// fsz := im.FontSize()
	// im.PushStyleVarFloat(im.StyleVar)
	if im.Button("HEY") {
		println("clicked")
	}
	im.Image(1, im.Vec2{float32(uidata.xx), float32(uidata.yy)})
	im.PopStyleColor()
	im.PopFont()
	// im.I
	// im.Image()
	im.End()
}
