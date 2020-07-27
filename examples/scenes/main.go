package main

import (
	"context"
	"time"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	osfs "github.com/gabstv/primen/io/os"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func main() {
	ebiten.SetRunnableOnUnfocused(true)
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:  1024,
		Height: 768,
		Scale:  ebiten.DeviceScaleFactor(),
		FS:     osfs.New("../shared"),
		Title:  "PRIMEN - Scenes",
		OnReady: func(e primen.Engine) {
			_, ch, _ := e.LoadScene("scene1")
			<-ch
		},
	})
	_ = engine.Run()
}

type SceneA struct {
	w         primen.World
	container io.Container
	engine    primen.Engine
	invalid   bool
}

func (s *SceneA) Engine() primen.Engine {
	return s.engine
}

func (*SceneA) Name() string {
	return "scene1"
}

func (s *SceneA) load() chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		w := s.engine.NewWorldWithDefaults(0)
		w.SetEnabled(false)
		s.w = w
		// load sprites
		s.container = io.NewContainer(context.Background(), s.engine.FS())
		_, done := s.container.LoadAll([]string{"particle2.png", "people.dat"})
		<-done
		atlas, err := s.container.GetAtlas("people.dat")
		if err != nil {
			panic(err)
		}
		root := primen.NewRootFnNode(s.w)
		root.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
			root.Transform().SetX(float64(s.engine.Width() / 2)).SetY(float64(s.engine.Height() / 2))
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.gotoscene2()
			}
		}
		sf := ebiten.DeviceScaleFactor()
		root.Transform().SetScale(sf*4, sf*4)
		spr := primen.NewChildAnimatedSpriteNode(root, primen.Layer0, 12, atlas.GetAnimation("boy"))
		spr.SpriteAnim().PlayClip("idle")
		ttl := primen.NewChildLabelNode(root, primen.Layer0)
		ttl.Label().SetText("Scene: "+s.Name()).SetOrigin(.5, .5)
		ttl.Transform().SetY(-30).SetScale(sf*.1, sf*.1)
		s.w.SetEnabled(true)
	}()
	return ch
}

func (s *SceneA) Unload() chan struct{} {
	ch := make(chan struct{})
	s.w.SetEnabled(false)
	s.engine.RemoveWorld(s.w)
	s.container.UnloadAll()
	s.w = nil
	//s.engine = nil
	s.container = nil
	close(ch)
	return ch
}

func (s *SceneA) gotoscene2() {
	if s.invalid {
		return
	}
	s.invalid = true
	go func() {
		_, sig, _ := s.engine.LoadScene("scene2")
		<-sig
		// s2.Start()
		s.engine.RunFn(func() {
			s.w.SetEnabled(false)
		})
		time.Sleep(time.Second)
		_ = s.Unload()
	}()
}

////////////////////////////////

type SceneB struct {
	w         primen.World
	container io.Container
	engine    primen.Engine
	invalid   bool
}

func (s *SceneB) Engine() primen.Engine {
	return s.engine
}

func (*SceneB) Name() string {
	return "scene2"
}

func (s *SceneB) load() chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		w := s.engine.NewWorldWithDefaults(0)
		w.SetEnabled(false)
		s.w = w
		// load sprites
		s.container = io.NewContainer(context.Background(), s.engine.FS())
		_, done := s.container.LoadAll([]string{"people.dat"})
		<-done
		atlas, err := s.container.GetAtlas("people.dat")
		if err != nil {
			panic(err)
		}
		root := primen.NewRootFnNode(s.w)
		root.Function().Update = func(ctx core.UpdateCtx, e ecs.Entity) {
			root.Transform().SetX(float64(s.engine.Width() / 2)).SetY(float64(s.engine.Height() / 2))
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.gotoscene1()
			}
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				s.Unload()
			}
		}
		sf := ebiten.DeviceScaleFactor()
		root.Transform().SetScale(sf*4, sf*4)
		spr := primen.NewChildAnimatedSpriteNode(root, primen.Layer0, 12, atlas.GetAnimation("girl"))
		spr.SpriteAnim().PlayClip("idle")
		ttl := primen.NewChildLabelNode(root, primen.Layer0)
		ttl.Label().SetText("Scene: "+s.Name()).SetOrigin(.5, .5)
		ttl.Transform().SetY(-30).SetScale(sf*.1, sf*.1)
	}()
	s.engine.RunFn(func() {
		// runs on main thread
		s.w.SetEnabled(true)
	})
	return ch
}

func (s *SceneB) Unload() chan struct{} {
	ch := make(chan struct{})
	s.w.SetEnabled(false)
	s.engine.RemoveWorld(s.w)
	s.container.UnloadAll()
	close(ch)
	s.w = nil
	//s.engine = nil
	s.container = nil
	return ch
}

func (s *SceneB) gotoscene1() {
	if s.invalid {
		return
	}
	s.invalid = true
	go func() {
		_, sig, _ := s.engine.LoadScene("scene1")
		<-sig
		s.engine.RunFn(func() {
			s.w.SetEnabled(false)
		})
		time.Sleep(time.Second)
		_ = s.Unload()
	}()
}

////////////////////////////////

func init() {
	primen.RegisterScene("scene1", func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &SceneA{engine: engine}
		ch := scn.load()
		return scn, ch
	})
	primen.RegisterScene("scene2", func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &SceneB{engine: engine}
		ch := scn.load()
		return scn, ch
	})
}
