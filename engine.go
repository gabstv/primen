package primen

import (
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io"
	osfs "github.com/gabstv/primen/io/os"
	"github.com/hajimehoshi/ebiten"
)

type StepInfo struct {
	l     sync.RWMutex
	lt    time.Time
	frame int64
}

func (i *StepInfo) Get() (lt time.Time, frame int64) {
	i.l.RLock()
	defer i.l.RUnlock()
	return i.lt, i.frame
}

func (i *StepInfo) GetFrame() (frame int64) {
	i.l.RLock()
	defer i.l.RUnlock()
	return i.frame
}

func (i *StepInfo) Set(lt time.Time, frame int64) {
	i.l.Lock()
	defer i.l.Unlock()
	i.lt = lt
	i.frame = frame
}

// Engine is what controls the ECS of primen.
type Engine struct {
	updateInfo   *StepInfo
	drawInfo     *StepInfo
	lock         sync.Mutex
	worlds       []worldContainer
	defaultWorld *core.GameWorld
	dmap         Dict
	options      EngineOptions
	f            io.Filesystem
	donech       chan struct{}
	once         sync.Once
	ready        func(e *Engine)
	ebilock      sync.RWMutex
	ebiOutsideW  int
	ebiOutsideH  int
	ebiLogicalW  int
	ebiLogicalH  int
	ebiScale     float64
	eventManager *core.EventManager
}

// NewEngineInput is the input data of NewEngine
type NewEngineInput struct {
	Width             int             // main window width
	Height            int             // main window height
	Scale             float64         // pixel scale (default: 1)
	TransparentScreen bool            // transparent screen
	Maximized         bool            // start window maximized
	Floating          bool            // always on top of all windows
	Fullscreen        bool            // start in fullscreen
	Resizable         bool            // is window resizable?
	FixedResolution   bool            // fixed logical screen resolution
	FixedWidth        int             // fixed logical screen resolution
	FixedHeight       int             // fixed logical screen resolution
	MaxResolution     bool            // set width/height to max resolution
	Title             string          // window title
	FS                io.Filesystem   // TODO: drop this
	OnReady           func(e *Engine) // function to run once the window is opened
}

// EngineOptions is used to setup Ebiten @ Engine.boot
type EngineOptions struct {
	Width               int
	Height              int
	Scale               float64
	Title               string
	IsFullscreen        bool
	IsResizable         bool
	IsMaxResolution     bool
	IsTransparentScreen bool
	IsFloating          bool
	IsMaximized         bool
	IsFixedResolution   bool
}

// Options will create a EngineOptions struct to be used in
// an *Engine
func (i *NewEngineInput) Options() EngineOptions {
	opt := EngineOptions{
		Width:               i.Width,
		Height:              i.Height,
		Scale:               i.Scale,
		Title:               i.Title,
		IsFullscreen:        i.Fullscreen,
		IsResizable:         i.Resizable,
		IsMaximized:         i.Maximized,
		IsMaxResolution:     i.MaxResolution,
		IsTransparentScreen: i.TransparentScreen,
		IsFloating:          i.Floating,
		IsFixedResolution:   i.FixedResolution,
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
			Width:             800,
			Height:            600,
			Scale:             1,
			Title:             "tau",
			FS:                osfs.New(fbase),
			FixedResolution:   false,
			Fullscreen:        false,
			Resizable:         false,
			MaxResolution:     false,
			TransparentScreen: false,
			Floating:          false,
		}
	} else {
		if v.Scale <= 0 {
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
	calcW, calcH := int(float64(v.Width)*v.Scale), int(float64(v.Height)*v.Scale)
	//if v.FixedResolution {
	//	calcW, calcH = v.Width, v.Height
	//}
	iw, ih := getLogicalSize(v.Width, v.Height, v.Scale, calcW, calcH, v.FixedResolution)
	if v.FixedWidth != 0 {
		iw = v.FixedWidth
	}
	if v.FixedHeight != 0 {
		ih = v.FixedHeight
	}

	e := &Engine{
		updateInfo:   &StepInfo{},
		drawInfo:     &StepInfo{},
		options:      v.Options(),
		f:            v.FS,
		donech:       make(chan struct{}),
		ready:        v.OnReady,
		ebiLogicalW:  iw,
		ebiLogicalH:  ih,
		ebiOutsideW:  v.Width,
		ebiOutsideH:  v.Height,
		ebiScale:     v.Scale,
		eventManager: &core.EventManager{},
	}

	// create the default world
	dw := core.NewWorld(e)
	// start default components and systems
	ecs.RegisterWorldDefaults(dw)

	e.worlds = []worldContainer{
		{
			priority: 0,
			world:    dw,
		},
	}
	e.defaultWorld = dw

	return e
}

// NewWorld adds a world to the engine.
// The priority is used to sort world execution, from hight to low.
func (e *Engine) NewWorld(priority int) *core.GameWorld {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.worlds == nil {
		e.worlds = make([]worldContainer, 0, 2)
	}
	ww := core.NewWorld(e)
	e.worlds = append(e.worlds, worldContainer{
		priority: priority,
		world:    ww,
	})
	// sort by priority
	sort.Sort(sortedWorldContainer(e.worlds))
	return ww
}

// RemoveWorld removes a *World
func (e *Engine) RemoveWorld(w *core.GameWorld) {
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
		return
	}
	// splice
	e.worlds = append(e.worlds[:wi], e.worlds[wi+1:]...)
	if w == e.defaultWorld {
		e.defaultWorld = nil
	}
}

// Default world
func (e *Engine) Default() *core.GameWorld {
	return e.defaultWorld
}

// Run boots up the game engine
func (e *Engine) Run() error {
	now := time.Now()
	e.drawInfo.Set(now, 0)
	e.updateInfo.Set(now, 0)

	ebiten.SetScreenTransparent(e.options.IsTransparentScreen)
	ebiten.SetFullscreen(e.options.IsFullscreen)
	ebiten.SetWindowResizable(e.options.IsResizable)
	ebiten.SetWindowFloating(e.options.IsFloating)
	if e.options.IsMaximized {
		ebiten.MaximizeWindow()
	}
	if e.options.IsMaxResolution {
		w, h := ebiten.WindowSize()
		if w != 0 && h != 0 {
			opt := e.options
			opt.Width = w
			opt.Height = h
			e.options = opt
		}
	}
	ebiten.SetWindowSize(e.options.Width, e.options.Height)
	ebiten.SetWindowTitle(e.options.Title)
	return ebiten.RunGame(e)
}

// Ready returns a channel that signals when the engine is ready
func (e *Engine) Ready() <-chan struct{} {
	return e.donech
}

// UpdateFrame returns the current frame. Use ctx.Frame() (more performant)
func (e *Engine) UpdateFrame() int64 {
	return e.updateInfo.GetFrame()
}

// DrawFrame returns the current frame. Use ctx.Frame() (more performant)
func (e *Engine) DrawFrame() int64 {
	return e.drawInfo.GetFrame()
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

// Width returns the logical width
func (e *Engine) Width() int {
	e.ebilock.RLock()
	defer e.ebilock.RUnlock()
	return e.ebiLogicalW
}

// Height returns the logical height
func (e *Engine) Height() int {
	e.ebilock.RLock()
	defer e.ebilock.RUnlock()
	return e.ebiLogicalH
}

// EBITEN Game interface

// Layout for ebiten.Game inteface
func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	e.ebilock.RLock()
	pow, poh := e.ebiOutsideW, e.ebiOutsideH
	piw, pih := e.ebiLogicalW, e.ebiLogicalH
	pscale := e.ebiScale
	pfixed := e.options.IsFixedResolution
	e.ebilock.RUnlock()
	niw, nih := getLogicalSize(outsideWidth, outsideHeight, pscale, piw, pih, pfixed)
	if outsideWidth == pow && outsideHeight == poh && piw == niw && pih == nih {
		return piw, pih
	}
	e.ebilock.Lock()
	defer e.ebilock.Unlock()
	e.ebiOutsideW = outsideWidth
	e.ebiOutsideH = outsideHeight
	e.ebiLogicalW = niw
	e.ebiLogicalH = nih
	return niw, nih
}

func (e *Engine) SetScreenScale(scale float64) {
	if scale <= 0 {
		return
	}
	e.ebiScale = scale
}

func (e *Engine) Update(screen *ebiten.Image) error {
	lastt, lastf := e.updateInfo.Get()
	now := time.Now()
	delta := now.Sub(lastt).Seconds()
	e.lock.Lock()
	worlds := e.worlds
	e.lock.Unlock()
	frame := lastf + 1
	e.updateInfo.Set(now, frame)

	e.once.Do(func() {
		close(e.donech)
		if e.ready != nil {
			e.ready(e)
		}
	})

	ctx := core.NewUpdateCtx(frame, delta, ebiten.CurrentTPS())

	for _, w := range worlds {
		w.world.EachSystem(func(s ecs.BaseSystem) bool {
			s.(core.System).UpdatePriority(ctx)
			return true
		})
		w.world.EachSystem(func(s ecs.BaseSystem) bool {
			s.(core.System).Update(ctx)
			return true
		})
	}
	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	lastt, lastf := e.drawInfo.Get()
	now := time.Now()
	delta := now.Sub(lastt).Seconds()
	//e.dmap.Set(TagDelta, delta) // set on update
	e.lock.Lock()
	worlds := e.worlds
	e.lock.Unlock()
	frame := lastf + 1
	e.drawInfo.Set(now, frame)

	//TODO: create DrawCtx
	ctx := core.NewDrawCtx(frame, delta, ebiten.CurrentTPS(), screen)

	for _, w := range worlds {
		w.world.EachSystem(func(s ecs.BaseSystem) bool {
			s.(core.System).DrawPriority(ctx)
			return true
		})
		w.world.EachSystem(func(s ecs.BaseSystem) bool {
			s.(core.System).Draw(ctx)
			return true
		})
	}
}

func getLogicalSize(outw, outh int, scale float64, inw, inh int, fixed bool) (w, h int) {
	if scale <= 0 {
		return outw, outh
	}
	if fixed {
		return inw, inh
		//return int(float64(inw) * scale), int(float64(inh) * scale)
	}
	return int(float64(outw) * scale), int(float64(outh) * scale)
}
