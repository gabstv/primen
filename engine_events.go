package primen

import "github.com/gabstv/primen/core"

func (e *engine) AddEventListener(eventName string, fn core.EventFn) core.EventID {
	return e.eventManager.Register(eventName, fn)
}
func (e *engine) RemoveEventListener(id core.EventID) bool {
	return e.eventManager.Deregister(id)
}
func (e *engine) DispatchEvent(eventName string, data interface{}) {
	e.eventManager.Dispatch(eventName, e, data)
}
