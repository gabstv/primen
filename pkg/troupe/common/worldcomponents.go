package common

import (
	"sync"

	"github.com/gabstv/ecs"
)

type WorldComponents struct {
	lock sync.RWMutex
	m    map[*ecs.World]*ecs.Component
}

func (wc *WorldComponents) Get(w *ecs.World) *ecs.Component {
	wc.lock.RLock()
	defer wc.lock.RUnlock()
	if wc.m == nil {
		return nil
	}
	return wc.m[w]
}

func (wc *WorldComponents) Set(w *ecs.World, c *ecs.Component) {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if wc.m == nil {
		wc.m = make(map[*ecs.World]*ecs.Component)
	}
	wc.m[w] = c
}
