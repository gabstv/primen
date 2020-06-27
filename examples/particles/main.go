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
		props.Duration = 1
		props.Source = []*ebiten.Image{pimg}
		props.InitColor = color.RGBA{
			R: 255,
			G: 100,
			B: 100,
			A: 255,
		}
		props.EndColor = color.RGBA{
			R: 0x00,
			G: 0xd6,
			B: 0xba,
			A: 0,
		}
		props.SetPositionRange(-3, 3, 0, 0)
		props.SetVelocityRange(-10, 10, -1, 1)
		props.XAccelVar0 = -50
		props.XAccelVar1 = 50
		props.YAccelVar0 = 0
		props.YAccelVar1 = 10
		props.InitScaleVar0 = .2
		props.InitScaleVar1 = 1
		pen.ParticleEmitter().SetProps(props).SetMaxParticles(200) //.SetX(50).SetY(50)
		em := pen.ParticleEmitter().EmissionProp()
		em.N0 = 2
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
	}
	{
		pen := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen.ParticleEmitter().Props()
		props.Duration = 1
		props.Source = []*ebiten.Image{pimg}
		props.InitColor = color.RGBA{
			R: 255,
			G: 100,
			B: 100,
			A: 255,
		}
		props.EndColor = color.RGBA{
			R: 60,
			G: 30,
			B: 255,
			A: 0,
		}
		props.XVelocityVar0 = -10
		props.XVelocityVar1 = 10
		props.XAccelVar0 = -50
		props.XAccelVar1 = 50
		props.YAccelVar0 = 0
		props.YAccelVar1 = 10
		pen.ParticleEmitter().SetProps(props).SetMaxParticles(200)
		pen.Transform().SetX(-100)
		em := pen.ParticleEmitter().EmissionProp()
		em.N0 = 2
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		pen.ParticleEmitter().SetCompositeMode(ebiten.CompositeModeLighter)
	}
	{
		pen3 := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen3.ParticleEmitter().Props()
		props.Duration = .25
		props.DurationVar1 = 2.5
		props.Source = []*ebiten.Image{pimg, pimg2, pimg3}
		props.InitColor = color.RGBA{
			R: 0xf1, // #f1c40f
			G: 0xc4,
			B: 0x0f,
			A: 255,
		}
		props.EndColor = color.RGBA{
			R: 0x8e, // #8e44ad
			G: 0x44,
			B: 0xad,
			A: 0,
		}
		props.YVelocity = -350
		props.YAccel = 180
		props.XVelocityVar0 = -14
		props.XVelocityVar1 = 14
		props.XAccelVar0 = -90
		props.XAccelVar1 = 90
		props.YAccelVar0 = -100
		props.YAccelVar1 = 100
		props.EndScaleVar0 = .2
		props.EndScaleVar1 = 2.9
		props.RotationVar0 = -1
		props.RotationVar1 = 1
		props.RotationAccelVar0, props.RotationAccelVar1 = -10, 10
		props.EndRotationAccelVar0, props.EndRotationAccelVar1 = -20, 20
		props.HueRotationSpeed = math.Pi / 2
		pen3.ParticleEmitter().SetProps(props).SetMaxParticles(2000)
		pen3.Transform().SetX(100)
		em := pen3.ParticleEmitter().EmissionProp()
		em.N0 = 5
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen3.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
	}
	penm := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
	{
		props := penm.ParticleEmitter().Props()
		props.Duration = .75
		props.DurationVar1 = 1.5
		props.Source = []*ebiten.Image{pimg, pimg2, pimg3}
		props.InitColor = color.RGBA{
			R: 0xfe, // #feca57
			G: 0xca,
			B: 0x57,
			A: 255,
		}
		props.EndColor = color.RGBA{
			R: 0xff,
			G: 0xff,
			B: 0xff,
			A: 0,
		}
		props.YVelocity = 0
		props.YAccel = 0
		props.XVelocityVar0 = -3
		props.XVelocityVar1 = 3
		props.YVelocityVar0 = -4
		props.YVelocityVar1 = 4
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
	}
	//tr.Transform().SetScale(1.8, 1.8)
	fn0 := primen.NewRootFnNode(w)
	var pgx, pgy float64
	var gx, gy float64
	var lx, ly float64
	dtt := 0.0
	dag1 := 0.0
	dag2 := math.Pi / 2
	dag3 := math.Pi / 4
	fn0.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		dag1 += math.Pi * ctx.DT() * 1.43123
		dag2 += math.Pi * ctx.DT() * 1.653
		dag3 += math.Pi * ctx.DT() * 0.3457823
		ss := (math.Cos(dag1) + math.Cos(dag2) + math.Sin(dag3)) / 3
		ssf := core.Lerpf(1, 1.6, (1+ss)/2)
		tr.Transform().SetScale(ssf, ssf)
		dtt += ctx.DT() / 4
		tr.Transform().SetAngle(dtt)
		tr.Transform().SetX(float64(ctx.Engine().Width()) / 2)
		tr.Transform().SetY(float64(ctx.Engine().Height()) / 2)
		//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		igx, igy := ebiten.CursorPosition()
		gx, gy = float64(igx), float64(igy)
		//fmt.Println(gx, gy)
		lx, ly, _ = core.GetTransformSystem(w).GlobalToLocal(gx, gy, tr.Entity())
		//fmt.Println(lx, ly)
		penm.Transform().SetX(lx).SetY(ly)
		//}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			ctx.Engine().SetScreenScale(ebiten.DeviceScaleFactor())
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			ctx.Engine().SetScreenScale(ebiten.DeviceScaleFactor() / 2)
		}
		props := penm.ParticleEmitter().Props()
		props.XVelocity = (gx - pgx) * 20
		props.YVelocity = (gy - pgy) * 20
		penm.ParticleEmitter().SetProps(props)
		pgx, pgy = gx, gy
	}
	fn0.Function().Draw = func(ctx core.DrawCtx, e ecs.Entity) {
		ebitenutil.DebugPrintAt(ctx.Renderer().Screen(), fmt.Sprintf("gx: %.4f\ngy: %.4f\nlx: %.4f\nly: %.4f\ntrx: %.4f\ntry: %.4f", gx, gy, lx, ly, tr.Transform().X(), tr.Transform().Y()), 10, 30)
	}
}
