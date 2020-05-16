package tau

import (
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/gabstv/ecs"
	"github.com/gabstv/tau/io"
	osfs "github.com/gabstv/tau/io/os"
	"github.com/hajimehoshi/ebiten"
)

// Engine is what controls the ECS of tau.
type Engine struct {
	lock         sync.Mutex
	lt           time.Time
	frame        int64
	worlds       []worldContainer
	defaultWorld *ecs.World
	dmap         Dict
	options      EngineOptions
	f            io.Filesystem
	donech       chan struct{}
	once         sync.Once
	ready        func(e *Engine)
}

// NewEngineInput is the input data of NewEngine
type NewEngineInput struct {
	Width   int
	Height  int
	Scale   float64
	Title   string
	FS      io.Filesystem
	OnReady func(e *Engine)
}

// EngineOptions is used to setup Ebiten @ Engine.boot
type EngineOptions struct {
	Width  int
	Height int
	Scale  float64
	Title  string
}

// Options will create a EngineOptions struct to be used in
// an *Engine
func (i *NewEngineInput) Options() EngineOptions {
	opt := EngineOptions{
		Width:  i.Width,
		Height: i.Height,
		Scale:  i.Scale,
		Title:  i.Title,
	}
	return opt
}

// NewEngine returns a new Engine
func NewEngine(v *NewEngineInput) *Engine {
	fbase := ""
	if len(os.Args) > 0 {
		fbase = path.Dir(os.Args[0])
	}
	if v == nil {
		v = &NewEngineInput{
			Width:  800,
			Height: 600,
			Scale:  1,
			Title:  "tau",
			FS:     osfs.New(fbase),
		}
	} else {
		if v.Scale < 1 {
			v.Scale = 1
		}
		if v.Width <= 0 {
			v.Width = 320
		}
		if v.Height <= 0 {
			v.Height = 240
		}
		if v.FS == nil {
			v.FS = osfs.New(fbase)
		}
	}
	// assign the default systems and controllers

	e := &Engine{
		options: v.Options(),
		f:       v.FS,
		donech:  make(chan struct{}),
		ready:   v.OnReady,
	}

	// create the default world
	dw := NewWorld(e)

	e.worlds = []worldContainer{
		worldContainer{
			priority: 0,
			world:    dw,
		},
	}
	e.defaultWorld = dw

	// start default components and systems
	startDefaults(e)

	return e
}

// AddWorld adds a world to the engine.
// The priority is used to sort world execution, from hight to low.
func (e *Engine) AddWorld(w *ecs.World, priority int) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.worlds == nil {
		e.worlds = make([]worldContainer, 0, 2)
	}
	e.worlds = append(e.worlds, worldContainer{
		priority: priority,
		world:    w,
	})
	// sort by priority
	sort.Sort(sortedWorldContainer(e.worlds))
}

// RemoveWorld removes a *World
func (e *Engine) RemoveWorld(w *ecs.World) bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	wi := -1
	for k, ww := range e.worlds {
		if ww.world == w {
			wi = k
			ww.world = nil
			break
		}
	}
	if wi == -1 {
		return false
	}
	// splice
	e.worlds = append(e.worlds[:wi], e.worlds[wi+1:]...)
	if w == e.defaultWorld {
		e.defaultWorld = nil
	}
	return true
}

// Default world
func (e *Engine) Default() *ecs.World {
	return e.defaultWorld
}

// Run boots up the game engine
func (e *Engine) Run() error {
	e.lock.Lock()
	width, height, scale, title := e.options.Width, e.options.Height, e.options.Scale, e.options.Title
	e.lt = time.Now()
	e.lock.Unlock()
	return ebiten.Run(e.loop, width, height, scale, title)
}

func (e *Engine) Ready() <-chan struct{} {
	return e.donech
}

func (e *Engine) loop(screen *ebiten.Image) error {
	e.lock.Lock()
	now := time.Now()
	ld := now.Sub(e.lt).Seconds()
	e.lt = now
	e.dmap.Set(TagDelta, ld)
	worlds := e.worlds
	e.frame++
	e.lock.Unlock()

	e.once.Do(func() {
		close(e.donech)
		if e.ready != nil {
			e.ready(e)
		}
	})

	for _, w := range worlds {
		w.world.Set("screen", screen)
		w.world.RunWithoutTag(WorldTagDraw, ld)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	for _, w := range worlds {
		w.world.Set("screen", screen)
		w.world.RunWithTag(WorldTagDraw, ld)
	}

	return nil
}

// Get an item from the global map
func (e *Engine) Get(key string) interface{} {
	return e.dmap.Get(key)
}

// Set an item to the global map
func (e *Engine) Set(key string, value interface{}) {
	e.dmap.Set(key, value)
}

// FS returns the filesystem
func (e *Engine) FS() io.Filesystem {
	return e.f
}
