package js

import (
	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
)

type sceneLoader interface {
	LastLoadedSceneJS() interface{}
	LoadSceneJS(name string) interface{}
}

type Scene interface {
	// Load() chan struct{}
	Unload() chan struct{}
}

type sceneLoadH interface {
	Scene() interface{}
	Ch() chan struct{}
	Err() error
}

func jsScene(runtime *goja.Runtime, e core.Engine, s Scene) *goja.Object {
	obj := runtime.NewObject()
	// obj.Set("load", func(call goja.FunctionCall) goja.Value {
	// 	ch := s.Load()
	// 	if fn, ok := goja.AssertFunction(call.Argument(0)); ok {
	// 		go func() {
	// 			<-ch
	// 			fn(call.This, obj)
	// 		}()
	// 	}
	// 	return goja.Null()
	// })
	obj.Set("unload", func(call goja.FunctionCall) goja.Value {
		ch := s.Unload()
		if fn, ok := goja.AssertFunction(call.Argument(0)); ok {
			go func() {
				<-ch
				fn(call.This, obj)
			}()
		}
		return goja.Null()
	})
	return obj
}

func sceneLast(e core.Engine, runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		si := e.(sceneLoader).LastLoadedSceneJS()
		if si == nil {
			return goja.Null()
		}
		return jsScene(runtime, e, si.(Scene))
	}
}

func scenesLoad(e core.Engine, runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		hh := e.(sceneLoader).LoadSceneJS(call.Argument(0).String()).(sceneLoadH)
		if hh.Err() != nil {
			return runtime.NewGoError(hh.Err())
		}
		so := jsScene(runtime, e, hh.Scene().(Scene))
		if fn, ok := goja.AssertFunction(call.Argument(1)); ok {
			go func() {
				<-hh.Ch()
				fn(call.This, so)
			}()
		}
		return so
	}
}
