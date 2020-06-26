package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var pimg *ebiten.Image
var pimg2 *ebiten.Image
var pimg3 *ebiten.Image

func main() {
	pimg, _, _ = ebitenutil.NewImageFromFile("../shared/particle.png", ebiten.FilterNearest)
	pimg2, _, _ = ebitenutil.NewImageFromFile("../shared/particle2.png", ebiten.FilterNearest)
	pimg3, _, _ = ebitenutil.NewImageFromFile("../shared/particle3.png", ebiten.FilterNearest)
	core.DebugDraw = true
	ebiten.SetRunnableOnUnfocused(true)
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     800,
		Height:    600,
		Resizable: true,
		OnReady:   ready,
		Scale:     ebiten.DeviceScaleFactor() / 2,
	})
	engine.Run()
}

func ready(engine primen.Engine) {
	engine.SetDebugTPS(true)
	w := engine.NewWorldWithDefaults(0)
	tr := primen.NewRootNode(w)
	coretr := tr.Transform()
	coretr.SetX(float64(engine.Width()) / 2).SetY(float64(engine.Height()) / 2)
	//
	{
		pen := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen.ParticleEmitter().Props()
		props.Dur = 1
		props.Source = []*ebiten.Image{pimg}
		props.Colorb = color.RGBA{
			R: 255,
			G: 100,
			B: 100,
			A: 255,
		}
		props.Colore = color.RGBA{
			R: 0x00,
			G: 0xd6,
			B: 0xba,
			A: 0,
		}
		//props.Vpx0 = -20
		//props.Vpx1 = 20
		props.Vvx0 = -10
		props.Vvx1 = 10
		props.Vax0 = -50
		props.Vax1 = 50
		props.Vay0 = 0
		props.Vay1 = 10
		props.InitScaleVar0 = .2
		props.InitScaleVar1 = 1
		pen.ParticleEmitter().SetProps(props).SetMaxParticles(200) //.SetX(50).SetY(50)
		em := pen.ParticleEmitter().EmissionProp()
		em.N0 = 2
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		pen.ParticleEmitter().SetEmissionParent(tr.Entity())
	}
	{
		pen := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen.ParticleEmitter().Props()
		props.Dur = 1
		props.Source = []*ebiten.Image{pimg}
		props.Colorb = color.RGBA{
			R: 255,
			G: 100,
			B: 100,
			A: 255,
		}
		props.Colore = color.RGBA{
			R: 60,
			G: 30,
			B: 255,
			A: 0,
		}
		props.Vvx0 = -10
		props.Vvx1 = 10
		props.Vax0 = -50
		props.Vax1 = 50
		props.Vay0 = 0
		props.Vay1 = 10
		pen.ParticleEmitter().SetProps(props).SetMaxParticles(200).SetX(-100)
		em := pen.ParticleEmitter().EmissionProp()
		em.N0 = 2
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		pen.ParticleEmitter().SetEmissionParent(tr.Entity())
		pen.ParticleEmitter().SetCompositeMode(ebiten.CompositeModeLighter)
	}
	{
		pen3 := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen3.ParticleEmitter().Props()
		props.Dur = .25
		props.Vdur1 = 2.5
		props.Source = []*ebiten.Image{pimg, pimg2, pimg3}
		props.Colorb = color.RGBA{
			R: 0xf1, // #f1c40f
			G: 0xc4,
			B: 0x0f,
			A: 255,
		}
		props.Colore = color.RGBA{
			R: 0x8e, // #8e44ad
			G: 0x44,
			B: 0xad,
			A: 0,
		}
		props.Vy = -350
		props.Ay = 180
		props.Vvx0 = -14
		props.Vvx1 = 14
		props.Vax0 = -90
		props.Vax1 = 90
		props.Vay0 = -100
		props.Vay1 = 100
		props.EndScaleVar0 = .2
		props.EndScaleVar1 = 2.9
		props.Vr0 = -1
		props.Vr1 = 1
		props.Vrab0, props.Vrab1 = -10, 10
		props.Vrae0, props.Vrae1 = -20, 20
		props.Hueshift = math.Pi / 2
		pen3.ParticleEmitter().SetProps(props).SetMaxParticles(2000).SetX(100)
		em := pen3.ParticleEmitter().EmissionProp()
		em.N0 = 5
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen3.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		pen3.ParticleEmitter().SetEmissionParent(tr.Entity())
	}
	penm := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
	{
		props := penm.ParticleEmitter().Props()
		props.Dur = .75
		props.Vdur1 = 1.5
		props.Source = []*ebiten.Image{pimg, pimg2, pimg3}
		props.Colorb = color.RGBA{
			R: 0xfe, // #feca57
			G: 0xca,
			B: 0x57,
			A: 255,
		}
		props.Colore = color.RGBA{
			R: 0xff,
			G: 0xff,
			B: 0xff,
			A: 0,
		}
		props.Vy = 0
		props.Ay = 0
		props.Vvx0 = -3
		props.Vvx1 = 3
		props.Vvy0 = -4
		props.Vvy1 = 4
		props.EndScaleVar0 = -1
		props.EndScaleVar1 = 1
		props.EndScale = 4
		props.InitScale = 0
		penm.ParticleEmitter().SetProps(props).SetMaxParticles(500)
		em := penm.ParticleEmitter().EmissionProp()
		em.N0 = 1
		em.N1 = 2
		em.T0 = .001
		em.T1 = .005
		penm.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		penm.ParticleEmitter().SetEmissionParent(tr.Entity())
	}
	//tr.Transform().SetScale(1.8, 1.8)
	fn0 := primen.NewRootFnNode(w)
	var pgx, pgy float64
	var gx, gy float64
	var lx, ly float64
	fn0.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		tr.Transform().SetX(float64(ctx.Engine().Width()) / 2)
		tr.Transform().SetY(float64(ctx.Engine().Height()) / 2)
		//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		igx, igy := ebiten.CursorPosition()
		gx, gy = float64(igx), float64(igy)
		//fmt.Println(gx, gy)
		lx, ly, _ = core.GetTransformSystem(w).GlobalToLocal(gx, gy, tr.Entity())
		//fmt.Println(lx, ly)
		penm.ParticleEmitter().SetX(lx).SetY(ly)
		//}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			ctx.Engine().SetScreenScale(ebiten.DeviceScaleFactor())
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			ctx.Engine().SetScreenScale(ebiten.DeviceScaleFactor() / 2)
		}
		props := penm.ParticleEmitter().Props()
		props.Vx = (gx - pgx) * 20
		props.Vy = (gy - pgy) * 20
		penm.ParticleEmitter().SetProps(props)
		pgx, pgy = gx, gy
	}
	fn0.Function().Draw = func(ctx core.DrawCtx, e ecs.Entity) {
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), fmt.Sprintf("gx: %.4f\ngy: %.4f\nlx: %.4f\nly: %.4f\ntrx: %.4f\ntry: %.4f", gx, gy, lx, ly, tr.Transform().X(), tr.Transform().Y()), 10, 30)
	}
}
