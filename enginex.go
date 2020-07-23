package primen

import (
	"context"
	"sort"

	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/geom"
	"github.com/gabstv/primen/io"
	"github.com/hajimehoshi/ebiten"
)

type Engine interface {
	core.Engine
	Ctx() context.Context
	Exit()
	FS() io.Filesystem
	LoadScene(name string) (scene Scene, sig chan struct{}, err error)
	RunFn(fn func())
	LastLoadedScene() Scene
	NewContainer() io.Container
	WaitAndGrabScreenImage() ScreenCopyRequest
	AddTempDrawFn(priority int, fn func(ctx core.DrawCtx) bool)
	NewDrawTarget(mask core.DrawMask, bounds geom.Rect, filter ebiten.Filter) core.DrawTargetID
	NewScreenOffsetDrawTarget(mask core.DrawMask) core.DrawTargetID
	NewProgrammableDrawTarget(input ProgrammableDrawTargetInput) core.DrawTargetID
	DrawTarget(id core.DrawTargetID) core.DrawTarget
	RemoveDrawTarget(id core.DrawTargetID) bool
}

type DrawCtx = core.DrawCtx

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

// NewContainer is a shorthand of io.NewContainer(engine.Ctx(), engine.FS())
func (e *engine) NewContainer() io.Container {
	return io.NewContainer(e.Ctx(), e.FS())
}

type drawFuncContainer struct {
	Priority int
	Func     func(ctx DrawCtx) bool
}

type sortedDrawFuncContainers []drawFuncContainer

func (a sortedDrawFuncContainers) Len() int           { return len(a) }
func (a sortedDrawFuncContainers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedDrawFuncContainers) Less(i, j int) bool { return a[i].Priority > a[j].Priority }

func (e *engine) AddTempDrawFn(priority int, fn func(ctx DrawCtx) bool) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.tempDrawFns = append(e.tempDrawFns, drawFuncContainer{
		Priority: priority,
		Func:     fn,
	})
	sort.Sort(sortedDrawFuncContainers(e.tempDrawFns))
}
