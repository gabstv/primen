package js

import (
	"errors"

	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
)

func eventDispatchFn(e core.Engine, runtime *goja.Runtime) func(name string, value goja.Value) error {
	return func(name string, value goja.Value) error {
		if name != "" {
			e.DispatchEvent(name, value.Export())
			return nil
		}
		return errors.New("event name must be a non empty string")
	}
}

func eventListenFn(e core.Engine, runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		arg0i := call.Argument(0).Export()
		if arg0i == nil {
			return runtime.ToValue(-1)
		}
		if callback, ok := goja.AssertFunction(call.Argument(1)); ok {
			if name, ok := arg0i.(string); ok {
				id := e.AddEventListener(name, func(eventName string, e core.Event) {
					callback(goja.Null(), runtime.ToValue(name), runtime.ToValue(e.Data))
				})
				return runtime.ToValue(int64(id))
			}
		}
		return runtime.ToValue(-1)
	}
}

func eventRemoveListenerFn(e core.Engine, runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		id := call.Argument(0).ToInteger()
		if id > 0 {
			e.RemoveEventListener(core.EventID(id))
		}
		return nil
	}
}
