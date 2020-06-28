package primen

import (
	"sync"
)

type Scene interface {
	Name() string
	Load(engine Engine) chan struct{}
	Unload() chan struct{}
	Start()
	Message(msg string)
}

type NewSceneFn func(engine Engine) Scene

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
	scene = e.sceneldrs[name](e)
	sig = scene.Load(e)
	return
}
