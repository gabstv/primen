package troupe

import (
	"sync"
)

type WorldComponents struct {
	lock sync.RWMutex
	m    map[Worlder]*Component
}

func (wc *WorldComponents) Get(w Worlder) *Component {
	wc.lock.RLock()
	defer wc.lock.RUnlock()
	if wc.m == nil {
		return nil
	}
	return wc.m[w]
}

func (wc *WorldComponents) Set(w Worlder, c *Component) {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if wc.m == nil {
		wc.m = make(map[Worlder]*Component)
	}
	wc.m[w] = c
}
