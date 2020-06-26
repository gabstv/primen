package main

import (
	"image/color"
	"math"

	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var pimg *ebiten.Image

func main() {
	pimg, _, _ = ebitenutil.NewImageFromFile("../shared/particle.png", ebiten.FilterDefault)
	core.DebugDraw = true
	ebiten.SetRunnableOnUnfocused(true)
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     800,
		Height:    600,
		Resizable: true,
		OnReady:   ready,
		Scale:     ebiten.DeviceScaleFactor(),
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
			R: 60,
			G: 30,
			B: 255,
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
		props.Vsx0 = .2
		props.Vsx1 = 1
		props.Vsy0 = .2
		props.Vsy1 = 1
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
		//props.Vpx0 = -20
		//props.Vpx1 = 20
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
		pen := primen.NewChildParticleEmitterNode(tr, primen.Layer0)
		props := pen.ParticleEmitter().Props()
		props.Dur = .5
		props.Vdur1 = 3
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
		//props.Vpx0 = -20
		//props.Vpx1 = 20
		props.Vy = -350
		props.Ay = 180
		props.Vvx0 = -10
		props.Vvx1 = 10
		props.Vax0 = -90
		props.Vax1 = 90
		props.Vay0 = -100
		props.Vay1 = 100
		props.Vesx0 = .2
		props.Vesx1 = 2.9
		props.Vesy0 = .2
		props.Vesy1 = 2.9
		props.Vr0 = -1
		props.Vr1 = 1
		props.Vrab0, props.Vrab1 = -10, 10
		props.Vrae0, props.Vrae1 = -20, 20
		props.Hueshift = math.Pi / 2
		pen.ParticleEmitter().SetProps(props).SetMaxParticles(2000).SetX(100)
		em := pen.ParticleEmitter().EmissionProp()
		em.N0 = 5
		em.N1 = 10
		em.T0 = .05
		em.T1 = .1
		pen.ParticleEmitter().SetEmissionProp(em).SetStrategy(core.SpawnReplace)
		pen.ParticleEmitter().SetEmissionParent(tr.Entity())
	}
	tr.Transform().SetScale(1.8, 1.8)
}
