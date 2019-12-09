package common

import (
	"sync"

	"github.com/gabstv/troupe/pkg/troupe/ecs"
)

type WorldComponents struct {
	lock sync.RWMutex
	m    map[ecs.Worlder]*ecs.Component
}

func (wc *WorldComponents) Get(w ecs.Worlder) *ecs.Component {
	wc.lock.RLock()
	defer wc.lock.RUnlock()
	if wc.m == nil {
		return nil
	}
	return wc.m[w]
}

func (wc *WorldComponents) Set(w ecs.Worlder, c *ecs.Component) {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if wc.m == nil {
		wc.m = make(map[ecs.Worlder]*ecs.Component)
	}
	wc.m[w] = c
}
