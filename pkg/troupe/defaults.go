package troupe

import (
	"sync"

	"github.com/gabstv/troupe/pkg/troupe/ecs"
)

// DefaultStarter is the function used to start a component or a system
type DefaultStarter func(e *Engine, w *ecs.World)

var (
	defaultCompStarters []DefaultStarter
	defaultSysStarters  []DefaultStarter
	defaultLock         sync.Mutex
)

// DefaultComp will add a component starter function
func DefaultComp(st DefaultStarter) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	if defaultCompStarters == nil {
		defaultCompStarters = make([]DefaultStarter, 0, 64)
	}
	defaultCompStarters = append(defaultCompStarters, st)
}

// DefaultSys will add a system starter function
func DefaultSys(st DefaultStarter) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	if defaultSysStarters == nil {
		defaultSysStarters = make([]DefaultStarter, 0, 64)
	}
	defaultSysStarters = append(defaultSysStarters, st)
}

func startDefaultSystems(e *Engine) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	defaultWorld := e.Default()
	if defaultWorld == nil {
		return
	}
	for _, starter := range defaultSysStarters {
		starter(e, defaultWorld)
	}
}

func startDefaultComponents(e *Engine) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	defaultWorld := e.Default()
	if defaultWorld == nil {
		return
	}
	for _, starter := range defaultCompStarters {
		starter(e, defaultWorld)
	}
}
