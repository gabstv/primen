package primen

import (
	"context"
	"sync"

	"github.com/gabstv/primen/io"
)

type Scene interface {
	Name() string
	Unload() chan struct{}
}

type NewSceneFn func(engine Engine) (Scene, chan struct{})

var registeredScenes map[string]NewSceneFn
var registeredScenesM sync.Mutex

// RegisterScene registers a scene
func RegisterScene(name string, fn NewSceneFn) {
	registeredScenesM.Lock()
	defer registeredScenesM.Unlock()
	if registeredScenes == nil {
		registeredScenes = make(map[string]NewSceneFn)
	}
	if _, ok := registeredScenes[name]; ok {
		panic(name + " scene already registered")
	}
	registeredScenes[name] = fn
}

func RegisteredScenes() map[string]NewSceneFn {
	registeredScenesM.Lock()
	defer registeredScenesM.Unlock()
	if registeredScenes == nil {
		return make(map[string]NewSceneFn)
	}
	c := make(map[string]NewSceneFn)
	for k, v := range registeredScenes {
		c[k] = v
	}
	return c
}

func (e *engine) loadScenes() {
	e.sceneldrs = RegisteredScenes()
}

func (e *engine) LoadScene(name string) (scene Scene, sig chan struct{}, err error) {
	return e.loadScene(name)
}

func (e *engine) loadScene(name string) (scene Scene, sig chan struct{}, err error) {
	if _, ok := e.sceneldrs[name]; !ok {
		//TODO: log error
		return nil, nil, ErrSceneNotFound
	}
	scene, sig = e.sceneldrs[name](e)
	e.lock.Lock()
	e.lastScn = scene
	e.lock.Unlock()
	// sig = scene.Load()
	return
}

func (e *engine) LastLoadedSceneJS() interface{} {
	return e.LastLoadedScene()
}

func (e *engine) LoadSceneJS(name string) interface{} {
	scn, ch, err := e.LoadScene(name)
	return &sceneLoaderHJS{
		scene: scn,
		sig:   ch,
		err:   err,
	}
}

type sceneLoaderHJS struct {
	scene Scene
	sig   chan struct{}
	err   error
}

func (sjs *sceneLoaderHJS) Scene() interface{} {
	return sjs.scene
}

func (sjs *sceneLoaderHJS) Ch() chan struct{} {
	return sjs.sig
}

func (sjs *sceneLoaderHJS) Err() error {
	return sjs.err
}

type SceneBase struct {
	Engine    Engine
	Container io.Container
}

func (s *SceneBase) Setup(engine Engine) {
	s.Engine = engine
	s.Container = io.NewContainer(context.Background(), engine.FS())
}

func (s *SceneBase) Destroy() {
	s.Container.UnloadAll()
}
