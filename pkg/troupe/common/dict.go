package common

import (
	"sync"
)

// Dict is a concurrency-safe map.
type Dict struct {
	lock  sync.RWMutex
	items map[string]interface{}
}

func (d *Dict) initialize() {
	if d.items == nil {
		d.items = make(map[string]interface{})
	}
}

// Get an item
func (d *Dict) Get(key string) interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.items == nil {
		return nil
	}
	return d.items[key]
}

// Set an item
func (d *Dict) Set(key string, value interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.initialize()
	d.items[key] = value
}

// Exists checks if the item exists
func (d *Dict) Exists(key string) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.items == nil {
		return false
	}
	_, ok := d.items[key]
	return ok
}
