package main

import (
	"context"
	_ "image/png"
	"math/rand"

	"github.com/gabstv/primen/examples/shared"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	"github.com/gabstv/primen/io/broccolifs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		Title:  "PRIMEN - AudioPlayer Example",
		OnReady: func(e primen.Engine) {
			c := io.NewContainer(context.Background(), e.FS())
			raw, err := c.GetAudioBytes("laser_shoot-3767.wav")
			if err != nil {
				panic(err)
			}
			w := e.NewWorldWithDefaults(0)
			ra := primen.NewRootAudioPlayerNode(w, core.NewAudioPlayerInput{
				RawAudio: raw,
				Panning:  true,
			})

			pem := pem1(w, c)
			pem2 := pem2(w, c)
			pem3 := pem3(w, c)

			lbl := primen.NewRootLabelNode(w, primen.Layer1)
			lbl.Label().SetText("Click to explode!")
			lbl.Transform().SetX(20).SetY(20)

			root := primen.NewRootFnNode(w)

			root.Function().Update = func(ctx core.UpdateCtx, en ecs.Entity) {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					xi, yi := ebiten.CursorPosition()
					pan := (float64(xi)/float64(e.Width()))*2 - 1
					ra.AudioPlayer().SetPan(pan)
					ra.AudioPlayer().Seek(0)
					ra.AudioPlayer().Play()
					pem.Transform().SetX(float64(xi)).SetY(float64(yi))
					e.RunFn(func() {})
					e.RunFn(func() {
						for i := 0; i < 6; i++ {
							pem3.ParticleEmitter().Emit(pem.Transform())
							pem.ParticleEmitter().Emit(pem.Transform())
							pem3.ParticleEmitter().Emit(pem.Transform())
							if rand.Float64() < .3 {
								p := pem.ParticleEmitter().Props()
								p2 := p
								p2.InitColor = primen.ColorFromHex("ffffff")
								p2.EndColor = primen.ColorFromHex("ffffff00")
								p2.InitScale = .1
								p2.InitScaleVar0 = 0
								p2.InitScaleVar1 = .1
								pem.ParticleEmitter().SetProps(p2)
								pem.ParticleEmitter().Emit(pem.Transform())
								pem.ParticleEmitter().Emit(pem.Transform())
								pem.ParticleEmitter().SetProps(p)
							}
						}
						pem2.ParticleEmitter().Emit(pem.Transform())
					})
				}
			}
		},
		FS: broccolifs.New(shared.Rx),
	})
	engine.Run()
}

func pem2(w primen.World, c io.Container) *primen.ParticleEmitterNode {
	pem := primen.NewRootParticleEmitterNode(w, 0)
	props := pem.ParticleEmitter().Props()
	props.Duration = .1
	props.InitScale = 3
	props.EndScale = 2.7
	props.InitColor = primen.ColorFromHex("ffffffff")
	props.EndColor = primen.ColorFromHex("ffffff00")
	props.OriginX = .5
	props.OriginY = .5
	rawimg, _ := c.GetImage("particle.png")
	img, _ := ebiten.NewImageFromImage(rawimg, ebiten.FilterNearest)
	props.Source = []*ebiten.Image{img}
	props.XVelocity = 0
	props.YVelocity = 0
	pem.ParticleEmitter().SetProps(props)
	pem.ParticleEmitter().SetCompositeMode(ebiten.CompositeModeLighter)
	ep := pem.ParticleEmitter().EmissionProp()
	ep.Enabled = false
	pem.ParticleEmitter().SetEmissionProp(ep)
	return pem
}

func pem1(w primen.World, c io.Container) *primen.ParticleEmitterNode {
	pem := primen.NewRootParticleEmitterNode(w, 0)
	props := pem.ParticleEmitter().Props()
	props.Duration = 1
	props.InitScale = .7
	props.EndScale = .3
	props.InitColor = primen.ColorFromHex("f68484ee")
	props.EndColor = primen.ColorFromHex("fffab200")
	props.OriginX = .5
	props.OriginY = .5
	rawimg, _ := c.GetImage("particle.png")
	img, _ := ebiten.NewImageFromImage(rawimg, ebiten.FilterNearest)
	props.Source = []*ebiten.Image{img}
	props.DurationVar0 = -.5
	props.DurationVar1 = 1
	props.XVelocityVar0 = -650
	props.XVelocityVar1 = 650
	props.YVelocityVar0 = -650
	props.YVelocityVar1 = 650
	props.XVelocity = 0
	props.YVelocity = 0
	props.InitScaleVar0 = -.2
	props.InitScaleVar1 = 0
	pem.ParticleEmitter().SetProps(props)
	pem.ParticleEmitter().SetCompositeMode(ebiten.CompositeModeLighter)
	ep := pem.ParticleEmitter().EmissionProp()
	ep.Enabled = false
	pem.ParticleEmitter().SetEmissionProp(ep)
	return pem
}

func pem3(w primen.World, c io.Container) *primen.ParticleEmitterNode {
	pem := primen.NewRootParticleEmitterNode(w, 0)
	props := pem.ParticleEmitter().Props()
	props.Duration = .3
	props.EndScale = .1
	props.InitColor = primen.ColorFromHex("fdffefee")
	props.EndColor = primen.ColorFromHex("fdffef00")
	props.OriginX = .5
	props.OriginY = .5
	rawimg, _ := c.GetImage("particle3.png")
	img, _ := ebiten.NewImageFromImage(rawimg, ebiten.FilterNearest)
	props.Source = []*ebiten.Image{img}
	props.DurationVar0 = 0
	props.DurationVar1 = .5
	props.XVelocityVar0 = -650
	props.XVelocityVar1 = 650
	props.YVelocityVar0 = -650
	props.YVelocityVar1 = 650
	props.XVelocity = 0
	props.YVelocity = 0
	props.InitScaleVar0 = -.8
	props.InitScaleVar1 = 0
	pem.ParticleEmitter().SetProps(props)
	pem.ParticleEmitter().SetCompositeMode(ebiten.CompositeModeLighter)
	ep := pem.ParticleEmitter().EmissionProp()
	ep.Enabled = false
	pem.ParticleEmitter().SetEmissionProp(ep)
	return pem
}
