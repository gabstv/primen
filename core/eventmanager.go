package core

import (
	"sync"
	"sync/atomic"
)

//FIXME: review

type EventID int64

type EventFn func(eventName string, e Event)

type EventManager struct {
	nextid     int64
	l          sync.RWMutex
	registered map[EventID]Observer
	list       map[string]map[EventID]Observer
}

func (m *EventManager) Register(eventName string, fn EventFn) EventID {
	m.l.Lock()
	defer m.l.Unlock()
	if m.registered == nil {
		m.registered = make(map[EventID]Observer)
	}
	if m.list == nil {
		m.list = make(map[string]map[EventID]Observer)
	}
	if m.list[eventName] == nil {
		m.list[eventName] = make(map[EventID]Observer)
	}
	id := EventID(atomic.AddInt64(&m.nextid, 1))
	o := &observer{
		name: eventName,
		fn:   fn,
	}
	m.list[eventName][id] = o
	m.registered[id] = o
	return id
}

func (m *EventManager) Deregister(id EventID) bool {
	ok := false
	var ob Observer
	m.l.RLock()
	ob, ok = m.registered[id]
	m.l.RUnlock()
	if !ok {
		return false
	}
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.registered, id)
	delete(m.list[ob.EventName()], id)
	return true
}

func (m *EventManager) Dispatch(eventName string, engine Engine, data interface{}) {
	m.l.RLock()
	if m.list == nil {
		// no listeners
		m.l.RUnlock()
		return
	}
	if m.list[eventName] == nil {
		// no listeners
		m.l.RUnlock()
		return
	}
	list := make([]Observer, 0, 8)
	for _, ob := range m.list[eventName] {
		list = append(list, ob)
	}
	m.l.RUnlock()
	for _, v := range list {
		go v.OnEvent(Event{
			Engine: engine,
			Data:   data,
		})
	}
}

type Observer interface {
	EventName() string
	OnEvent(e Event)
}

type observer struct {
	name string
	fn   EventFn
}

func (o *observer) EventName() string {
	return o.name
}

func (o *observer) OnEvent(e Event) {
	if o.fn != nil {
		o.fn(o.name, e)
	}
}

type Event struct {
	Engine Engine
	Data   interface{}
}
