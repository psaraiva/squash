//go:build js && wasm

package web

import (
	"syscall/js"
)

type Canvas struct {
	ctx JSContext
	w   float64
	h   float64
}

func NewRenderer() *Canvas {
	doc := js.Global().Get("document")
	canvasElement := doc.Call("getElementById", "gameCanvas")
	return NewCanvas(canvasElement)
}

func NewCanvas(canvasElement js.Value) *Canvas {
	ctx := canvasElement.Call("getContext", "2d")
	return &Canvas{
		ctx: NewJSContext(ctx),
		w:   canvasElement.Get("width").Float(),
		h:   canvasElement.Get("height").Float(),
	}
}

func (c *Canvas) Clear() {
	c.ctx.Set("fillStyle", "black")
	c.ctx.Call("fillRect", 0, 0, c.w, c.h)
}

func (c *Canvas) DrawBall(x, y, size float64) {
	c.ctx.Set("fillStyle", "white")
	c.ctx.Call("fillRect", x, y, size, size)
}

func (c *Canvas) DrawPaddle(x, y, w, h float64) {
	c.ctx.Set("fillStyle", "white")
	c.ctx.Call("fillRect", x, y, w, h)
}

func (c *Canvas) DrawText(text string, x, y float64) {
	c.ctx.Set("fillStyle", "white")
	c.ctx.Set("font", "20px Arial")
	c.ctx.Call("fillText", text, x, y)
}

func (c *Canvas) DrawDebugText(text string, x, y float64) {
	c.ctx.Set("font", "12px monospace")
	c.ctx.Set("fillStyle", "#FFFF00")
	c.ctx.Call("fillText", text, x, y)
}

func (c *Canvas) MeasureText(text string) float64 {
	metrics := c.ctx.Call("measureText", text)
	return metrics.Get("width").Float()
}
