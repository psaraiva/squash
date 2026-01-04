//go:build js && wasm

package web

import "syscall/js"

// JSContext interface abstracts syscall/js.Value operations
type JSContext interface {
	Call(method string, args ...interface{}) JSContext
	Float() float64
	Get(key string) JSContext
	Set(key string, value interface{})
}

// RealJSContext wraps js.Value to implement JSContext interface
type RealJSContext struct {
	value js.Value
}

// NewJSContext creates a new JSContext wrapper from js.Value
func NewJSContext(value js.Value) JSContext {
	return &RealJSContext{value: value}
}

func (r *RealJSContext) Set(key string, value interface{}) {
	r.value.Set(key, value)
}

func (r *RealJSContext) Call(method string, args ...interface{}) JSContext {
	result := r.value.Call(method, args...)
	return &RealJSContext{value: result}
}

func (r *RealJSContext) Get(key string) JSContext {
	result := r.value.Get(key)
	return &RealJSContext{value: result}
}

func (r *RealJSContext) Float() float64 {
	return r.value.Float()
}
