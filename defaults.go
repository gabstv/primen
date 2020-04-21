package tau

import (
	"sync"
)

var (
	defaultStartersMap map[string]struct{}
	defaultCSStarters  []ComponentSystem
	defaultLock        sync.Mutex
)

func RegisterComponentSystem(cs ComponentSystem) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	if defaultStartersMap == nil {
		defaultStartersMap = make(map[string]struct{})
		defaultCSStarters = make([]ComponentSystem, 0, 128)
	}
	if _, ok := defaultStartersMap[cs.SystemName()]; ok {
		panic("system " + cs.SystemName() + " already registered")
	}
	defaultStartersMap[cs.SystemName()] = struct{}{}
	defaultCSStarters = append(defaultCSStarters, cs)
}

func startDefaults(e *Engine) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	defaultWorld := e.Default()
	if defaultWorld == nil {
		return
	}
	for _, starter := range defaultCSStarters {
		_ = starter.Components(defaultWorld)
	}
	for _, starter := range defaultCSStarters {
		SetupSystem(defaultWorld, starter)
	}
}
