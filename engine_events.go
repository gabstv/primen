package tau

import "github.com/gabstv/tau/core"

func (e *Engine) AddEventListener(eventName string, fn core.EventFn) core.EventID {
	return e.eventManager.Register(eventName, fn)
}
func (e *Engine) RemoveEventListener(id core.EventID) bool {
	return e.eventManager.Deregister(id)
}
func (e *Engine) DispatchEvent(eventName string, data interface{}) {
	e.eventManager.Dispatch(eventName, e, data)
}
