package js

import (
	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
)

// Stdlib loads the Primen standard library into a js runtime
func Stdlib(e core.Engine, r *goja.Runtime) {
	events := r.NewObject()
	events.Set("dispatch", eventDispatchFn(e, r))
	events.Set("add_listener", eventListenFn(e, r))
	events.Set("remove_listener", eventRemoveListenerFn(e, r))
	r.Set("$events", events)

	scenes := r.NewObject()
	scenes.Set("last", sceneLast(e, r))
	r.Set("$scenes", scenes)
}
