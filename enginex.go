package primen

import (
	"context"

	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
)

type Engine interface {
	core.Engine
	Ctx() context.Context
	Exit()
	FS() io.Filesystem
	LoadScene(name string) (scene Scene, sig chan struct{}, err error)
	RunFn(fn func())
	LastLoadedScene() Scene
}

type World = core.World

// RunFn runs a function on the main thread (and before the ECS), so it's safe
// to use it for non thread safe stuff.
//
// The engine will only run one fn per frame.
func (e *engine) RunFn(fn func()) {
	e.runfns <- fn
}

func (e *engine) Exit() {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.exits = true
}

func (e *engine) LastLoadedScene() Scene {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.lastScn
}
