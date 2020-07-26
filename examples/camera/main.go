package main

import (
	_ "image/png"
	"math"
	"sync"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/core/ui/imgui"
	"github.com/gabstv/primen/dom"
	"github.com/gabstv/primen/examples/camera/res"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
	//"github.com/pkg/profile"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  800,
		Height: 600,
		FS:     res.FS(),
		OnReady: func(e primen.Engine) {
			e.LoadScene("main_menu")
		},
	})
	imgui.Setup(engine)
	engine.SetDebugFPS(true)
	engine.SetDebugTPS(true)
	engine.Run()
}

type MainMenuScene struct {
	engine primen.Engine
	ui     imgui.UID
}

var _ primen.Scene = (*MainMenuScene)(nil)

func (*MainMenuScene) Name() string {
	return "main_menu"
}

func (s *MainMenuScene) Unload() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	imgui.RemoveUI(s.ui)
	return ch
}

// PrevSceneCh implements AutoScene
func (s *MainMenuScene) PrevSceneCh(ch <-chan struct{}) {}

func (s *MainMenuScene) setup() chan struct{} {
	ch := make(chan struct{})
	c := s.engine.NewContainer()
	go func() {
		defer close(ch)
		_, done := c.LoadAll([]string{
			"public/mainmenu.xml",
		})
		<-done
		// run on main thread
		s.engine.RunFn(func() {
			node, _ := c.GetXMLDOM("public/mainmenu.xml")
			s.ui = imgui.AddUI(node.(dom.ElementNode))
		})
	}()
	return ch
}

type SinglePlayerScene struct {
	engine  primen.Engine
	ui      imgui.UID
	bgimage *ebiten.Image
	fgimage *ebiten.Image
	arimage *ebiten.Image
	w       primen.World
	s00     core.DrawTargetID
}

var _ primen.Scene = (*SinglePlayerScene)(nil)

func (*SinglePlayerScene) Name() string {
	return "single"
}

func (s *SinglePlayerScene) Unload() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	imgui.RemoveUI(s.ui)
	s.engine.RemoveWorld(s.w)
	s.engine.RemoveDrawTarget(s.s00)
	s.w = nil
	return ch
}

// PrevSceneCh implements AutoScene
func (s *SinglePlayerScene) PrevSceneCh(ch <-chan struct{}) {}

func (s *SinglePlayerScene) setup() chan struct{} {
	ch := make(chan struct{})
	c := s.engine.NewContainer()
	go func() {
		defer close(ch)
		_, done := c.LoadAll([]string{
			"public/single.xml",
			"public/cambg.png",
			"public/target.png",
			"public/arrow.png",
		})
		<-done
		var wg sync.WaitGroup
		wg.Add(2)
		// run on main thread
		s.engine.RunFn(func() {
			defer wg.Done()
			node, _ := c.GetXMLDOM("public/single.xml")
			s.ui = imgui.AddUI(node.(dom.ElementNode))
		})
		s.engine.RunFn(func() {
			defer wg.Done()
			im0, _ := c.GetImage("public/cambg.png")
			s.bgimage, _ = ebiten.NewImageFromImage(im0, ebiten.FilterNearest)
			im1, _ := c.GetImage("public/target.png")
			s.fgimage, _ = ebiten.NewImageFromImage(im1, ebiten.FilterNearest)
			im2, _ := c.GetImage("public/arrow.png")
			s.arimage, _ = ebiten.NewImageFromImage(im2, ebiten.FilterNearest)
			s.setupScene()
		})
		wg.Wait()
	}()
	return ch
}

func (s *SinglePlayerScene) setupScene() {
	s.s00 = s.engine.NewScreenOffsetDrawTarget(core.DrawMaskDefault)
	bgimage := s.bgimage
	fgimage := s.fgimage
	arimage := s.arimage
	e := s.engine
	w := e.NewWorldWithDefaults(0)
	s.w = w
	wtrn := primen.NewRootNode(w)
	ctr := primen.NewChildNode(wtrn)
	components.SetCameraComponentData(w, ctr.Entity(), components.NewCamera(s.s00))
	cam := components.GetCameraComponentData(w, ctr.Entity())
	cam.SetViewRect(e.SizeVec())
	components.SetFollowTransformComponentData(w, ctr.Entity(), components.FollowTransform{})
	c := components.GetFollowTransformComponentData(w, ctr.Entity())
	//c.SetDeadZone(geom.Rect{geom.ZV, geom.Vec{100, 100}}.At(e.SizeVec().Scaled(.5).Sub(geom.Vec{50, 50})))
	tiled := primen.NewChildTileSetNode(wtrn, primen.Layer0, []*ebiten.Image{bgimage}, 32, 32, 256, 256, make([]int, 32*32))
	_ = tiled
	dude := primen.NewChildSpriteNode(wtrn, primen.Layer1)
	dude.Sprite().SetImage(fgimage).SetOrigin(.5, .5)
	dude.Transform().SetPos(e.SizeVec().Scaled(.5))
	ctr.Transform().SetPos(e.SizeVec().Scaled(.5))
	ar1 := primen.NewChildSpriteNode(dude, primen.Layer2)
	ar1.Transform().SetScale(3, 3)
	ar1.Sprite().SetImage(arimage).SetOrigin(.5, .5)
	c.SetTarget(dude.Entity())
	c.SetBounds(geom.Rect{
		Min: e.SizeVec().Scaled(.5),
		Max: geom.Vec{2048, 2048},
	})
	c.SetDeadZone(geom.Vec{100, 40})
	//
	fnn := primen.NewRootFnNode(w)
	fnn.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		dtr := dude.Transform()
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			dtr.SetX(dtr.X() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			dtr.SetX(dtr.X() - (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			dtr.SetY(dtr.Y() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			dtr.SetY(dtr.Y() - (ctx.DT() * 200))
		}
	}
}

//

type DoublePlayerScene struct {
	engine   primen.Engine
	ui       imgui.UID
	bgimage  *ebiten.Image
	fgimage  *ebiten.Image
	arimage  *ebiten.Image
	particle *ebiten.Image
	w        primen.World
	s00      core.DrawTargetID
	s01      core.DrawTargetID
}

var _ primen.Scene = (*DoublePlayerScene)(nil)

func (*DoublePlayerScene) Name() string {
	return "double"
}

func (s *DoublePlayerScene) Unload() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	imgui.RemoveUI(s.ui)
	s.engine.RemoveWorld(s.w)
	s.engine.RemoveDrawTarget(s.s00)
	s.engine.RemoveDrawTarget(s.s01)
	s.w = nil
	return ch
}

// PrevSceneCh implements AutoScene
func (s *DoublePlayerScene) PrevSceneCh(ch <-chan struct{}) {}

func (s *DoublePlayerScene) setup() chan struct{} {
	ch := make(chan struct{})
	c := s.engine.NewContainer()
	go func() {
		defer close(ch)
		_, done := c.LoadAll([]string{
			"public/double.xml",
			"public/cambg.png",
			"public/target.png",
			"public/arrow.png",
			"public/particle3.png",
		})
		<-done
		var wg sync.WaitGroup
		wg.Add(2)
		// run on main thread
		s.engine.RunFn(func() {
			defer wg.Done()
			node, _ := c.GetXMLDOM("public/double.xml")
			s.ui = imgui.AddUI(node.(dom.ElementNode))
		})
		s.engine.RunFn(func() {
			defer wg.Done()
			im0, _ := c.GetImage("public/cambg.png")
			s.bgimage, _ = ebiten.NewImageFromImage(im0, ebiten.FilterNearest)
			im1, _ := c.GetImage("public/target.png")
			s.fgimage, _ = ebiten.NewImageFromImage(im1, ebiten.FilterNearest)
			im2, _ := c.GetImage("public/arrow.png")
			s.arimage, _ = ebiten.NewImageFromImage(im2, ebiten.FilterNearest)
			im3, _ := c.GetImage("public/particle3.png")
			s.particle, _ = ebiten.NewImageFromImage(im3, ebiten.FilterNearest)
			s.setupScene()
		})
		wg.Wait()
	}()
	return ch
}

func (s *DoublePlayerScene) setupScene() {
	bgimage := s.bgimage
	fgimage := s.fgimage
	arimage := s.arimage
	e := s.engine

	w := e.NewWorldWithDefaults(0)
	s.w = w

	s.s00 = e.NewDrawTarget(core.DrawMaskDefault, geom.Rect{Max: geom.Vec{.495, 1}}, ebiten.FilterDefault)
	s.s01 = e.NewDrawTarget(core.DrawMaskDefault, geom.Rect{Min: geom.Vec{0.505, 0}, Max: geom.Vec{1, 1}}, ebiten.FilterDefault)

	wtrn := primen.NewRootNode(w)

	tiled := primen.NewChildTileSetNode(wtrn, primen.Layer0, []*ebiten.Image{bgimage}, 16, 16, 256, 256, make([]int, 16*16))
	_ = tiled

	var player1 *primen.SpriteNode
	var player2 *primen.SpriteNode
	var ar1, ar2 *primen.SpriteNode
	var psys1 *primen.ParticleEmitterNode
	var psys2 *primen.ParticleEmitterNode
	// player 1
	{
		player1 = primen.NewChildSpriteNode(wtrn, primen.Layer1)
		player1.Sprite().SetImage(fgimage).SetOrigin(.5, .5)
		player1.Transform().SetPos(e.SizeVec().Scaled(.5))
		ar1 = primen.NewChildSpriteNode(player1, primen.Layer2)
		ar1.Transform().SetScale(3, 3)
		ar1.Sprite().SetImage(arimage).SetOrigin(.5, .5)
		psys1 = particleSys(wtrn, s.particle)
	}
	// player 2
	{
		player2 = primen.NewChildSpriteNode(wtrn, primen.Layer1)
		player2.Sprite().SetImage(fgimage).SetOrigin(.5, .5).RotateHue(math.Pi)
		player2.Transform().SetPos(e.SizeVec().Scaled(.5).Add(geom.Vec{0, 50}))
		ar2 = primen.NewChildSpriteNode(player2, primen.Layer2)
		ar2.Transform().SetScale(3, 3)
		ar2.Sprite().SetImage(arimage).SetOrigin(.5, .5)
		psys2 = particleSys(wtrn, s.particle)
	}
	// camera 1
	{
		camtr := primen.NewChildNode(wtrn)
		components.SetCameraComponentData(w, camtr.Entity(), components.NewCamera(s.s00))
		cam := components.GetCameraComponentData(w, camtr.Entity())
		cam.SetViewRect(e.SizeVec().ScaledXY(.5, 1))
		components.SetFollowTransformComponentData(w, camtr.Entity(), components.FollowTransform{})
		ftr := components.GetFollowTransformComponentData(w, camtr.Entity())
		ftr.SetTarget(player1.Entity())
		ftr.SetBounds(geom.Rect{
			Min: e.SizeVec().Scaled(.5).ScaledXY(.5, 1),
			Max: geom.Vec{2048, 2048},
		})
		ftr.SetDeadZone(geom.Vec{80, 40})
		camtr.Transform().SetPos(player1.Transform().Pos())
	}
	// camera 2
	{
		camtr := primen.NewChildNode(wtrn)
		components.SetCameraComponentData(w, camtr.Entity(), components.NewCamera(s.s01))
		cam := components.GetCameraComponentData(w, camtr.Entity())
		cam.SetViewRect(e.SizeVec().ScaledXY(.5, 1))
		components.SetFollowTransformComponentData(w, camtr.Entity(), components.FollowTransform{})
		ftr := components.GetFollowTransformComponentData(w, camtr.Entity())
		ftr.SetTarget(player2.Entity())
		ftr.SetBounds(geom.Rect{
			Min: e.SizeVec().Scaled(.5).ScaledXY(.5, 1),
			Max: geom.Vec{2048, 2048},
		})
		ftr.SetDeadZone(geom.Vec{80, 40})
		camtr.Transform().SetPos(player2.Transform().Pos())
	}
	//
	fnn := primen.NewRootFnNode(w)
	fnn.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
		a2a := player1.Transform().Pos().Sub(player2.Transform().Pos()).Angle()
		a1a := player2.Transform().Pos().Sub(player1.Transform().Pos()).Angle()
		ar1.Transform().SetAngle(a1a)
		ar2.Transform().SetAngle(a2a)
		dtr := player2.Transform()
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			dtr.SetX(dtr.X() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			dtr.SetX(dtr.X() - (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			dtr.SetY(dtr.Y() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			dtr.SetY(dtr.Y() - (ctx.DT() * 200))
		}
		dtr2 := player1.Transform()
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			dtr2.SetX(dtr2.X() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			dtr2.SetX(dtr2.X() - (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			dtr2.SetY(dtr2.Y() + (ctx.DT() * 200))
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			dtr2.SetY(dtr2.Y() - (ctx.DT() * 200))
		}
		if dtr.X() < 20 {
			dtr.SetX(20)
		}
		if dtr2.X() < 20 {
			dtr2.SetX(20)
		}
		if dtr.Y() < 10 {
			dtr.SetY(10)
		}
		if dtr2.Y() < 10 {
			dtr2.SetY(10)
		}
		if dtr.X() > 2038 {
			dtr.SetX(2038)
		}
		if dtr2.X() > 2038 {
			dtr2.SetX(2038)
		}
		if dtr.Y() > 2038 {
			dtr.SetY(2038)
		}
		if dtr2.Y() > 2038 {
			dtr2.SetY(2038)
		}
		psys1.Transform().SetPos(dtr.Pos())
		psys2.Transform().SetPos(dtr2.Pos())
	}
}

func particleSys(parent primen.ObjectContainer, img *ebiten.Image) *primen.ParticleEmitterNode {
	p := primen.NewChildParticleEmitterNode(parent, primen.Layer0)
	pe := p.ParticleEmitter()
	pe.SetCompositeMode(ebiten.CompositeModeLighter)
	props := pe.Props()
	props.DurationVar1 = .3
	props.Source = []*ebiten.Image{img}
	props.EndColor = primen.ColorFromHex("#ff000000")
	props.InitScale = 0.2
	props.InitScaleVar0 = 0.1
	props.InitScaleVar1 = 2
	props.InitColor = primen.ColorFromHex("#ffffff33")
	props.RotationAccelVar0 = -0.1
	props.RotationAccelVar1 = 0.1
	props.XVelocity = 0
	props.YVelocity = 0
	props.XVelocityVar0 = -10
	props.XVelocityVar1 = 10
	props.YVelocityVar0 = -10
	props.YVelocityVar1 = 10
	pe.SetProps(props)
	eprop := pe.EmissionProp()
	eprop.N0 = 2
	eprop.N1 = 6
	eprop.T0 = 1 / 60
	eprop.T1 = 2 / 60
	pe.SetEmissionProp(eprop)
	pe.SetMaxParticles(512)
	return p
}

//

func init() {
	primen.RegisterScene((*MainMenuScene)(nil).Name(), func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &MainMenuScene{
			engine: engine,
		}
		ch := scn.setup()
		return scn, ch
	})
	primen.RegisterScene((*SinglePlayerScene)(nil).Name(), func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &SinglePlayerScene{
			engine: engine,
		}
		ch := scn.setup()
		return scn, ch
	})
	primen.RegisterScene((*DoublePlayerScene)(nil).Name(), func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &DoublePlayerScene{
			engine: engine,
		}
		ch := scn.setup()
		return scn, ch
	})
}
