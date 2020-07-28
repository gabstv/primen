package main

import (
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	osfs "github.com/gabstv/primen/io/os"
	"github.com/gabstv/primen/modules/imgui"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	engine := primen.NewEngine(&primen.NewEngineInput{
		Width:     800,
		Height:    600,
		Resizable: true,
		Title:     "UI Scenes",
		Scene:     "default",
		FS:        osfs.New("assets"),
		OnReady: func(e primen.Engine) {
			imgui.Setup(e)
		},
	})
	engine.SetDebugTPS(true)
	engine.Run()
}

type DefaultScene struct {
	engine    primen.Engine
	container io.Container
	uiID      imgui.UID
}

func (s *DefaultScene) init() chan struct{} {
	s.container = s.engine.NewContainer()
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		_, donech := s.container.LoadAll([]string{"default.xml"})
		<-donech
		doms, err := s.container.GetXMLDOM("default.xml")
		if err != nil {
			panic(err)
		}
		s.engine.RunFn(func() {
			s.uiID = imgui.AddUI(doms)
		})
	}()
	return nil
}

func (s *DefaultScene) destroy() {}

func (s *DefaultScene) Name() string {
	return "default"
}

func (s *DefaultScene) Unload() chan struct{} {
	ch := make(chan struct{})
	go func() {
		imgRequest := s.engine.WaitAndGrabScreenImage()
		<-imgRequest.Done()
		img := imgRequest.ScreenCopy()
		opt := &ebiten.DrawImageOptions{}
		t := 1.0
		s.engine.AddTempDrawFn(0, func(ctx primen.DrawCtx) bool {
			opt.ColorM.Reset()
			opt.ColorM.Scale(1, 1, 1, t)
			ctx.Renderer().DrawImage(img, opt, core.DrawMaskAll)
			t -= 1.0 / 120.0
			if t <= 0 {
				t = 0
				close(ch)
				return false
			}
			return true
		})
		s.engine.RunFn(func() {
			imgui.RemoveUI(s.uiID)
			s.container.UnloadAll()
		})
	}()
	return ch
}

// PrevSceneCh implements AutoScene
func (s *DefaultScene) PrevSceneCh(ch <-chan struct{}) {}

type GameScene struct {
	engine    primen.Engine
	container io.Container
	uiID      imgui.UID
}

func (s *GameScene) init() chan struct{} {
	s.container = s.engine.NewContainer()
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		_, donech := s.container.LoadAll([]string{"game.xml"})
		<-donech
		doms, err := s.container.GetXMLDOM("game.xml")
		if err != nil {
			panic(err)
		}
		s.engine.RunFn(func() {
			s.uiID = imgui.AddUI(doms)
		})
	}()
	return nil
}

func (s *GameScene) Name() string {
	return "game"
}

func (s *GameScene) Unload() chan struct{} {
	ch := make(chan struct{})
	go func() {
		imgRequest := s.engine.WaitAndGrabScreenImage()
		<-imgRequest.Done()
		img := imgRequest.ScreenCopy()
		opt := &ebiten.DrawImageOptions{}
		t := 1.0
		s.engine.AddTempDrawFn(0, func(ctx primen.DrawCtx) bool {
			opt.ColorM.Reset()
			opt.ColorM.Scale(1, 1, 1, t)
			ctx.Renderer().DrawImage(img, opt, core.DrawMaskAll)
			t -= 1.0 / 120.0
			if t <= 0 {
				t = 0
				close(ch)
				return false
			}
			return true
		})
		s.engine.RunFn(func() {
			imgui.RemoveUI(s.uiID)
			s.container.UnloadAll()
		})
	}()
	return ch
}

// PrevSceneCh implements AutoScene
func (s *GameScene) PrevSceneCh(ch <-chan struct{}) {}

func init() {
	primen.RegisterScene("default", func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &DefaultScene{
			engine: engine,
		}
		ch := scn.init()
		return scn, ch
	})
	primen.RegisterScene("game", func(engine primen.Engine) (primen.Scene, chan struct{}) {
		scn := &GameScene{
			engine: engine,
		}
		ch := scn.init()
		return scn, ch
	})
}
