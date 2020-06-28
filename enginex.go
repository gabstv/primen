package primen

import (
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
)

type Engine interface {
	core.Engine
	FS() io.Filesystem
	LoadScene(name string) (scene Scene, sig chan struct{}, err error)
	RunFn(fn func())
}

type World = core.World

// RunFn runs a function on the main thread (and before the ECS), so it's safe
// to use it for non thread safe stuff.
//
// The engine will only run one fn per frame.
func (e *engine) RunFn(fn func()) {
	e.runfns <- fn
}
