package wasm

import (
	"syscall/js"

	"github.com/psaraiva/squash/internal/app"
)

func SetupMouseHandlers(squash *app.Squash, canvas js.Value, cfg app.Config) {
	// Reset/Start
	canvas.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) any {
		button := args[0].Get("button").Int()

		// left button
		if button == 0 && (squash.State == app.StateMenu || squash.State == app.StateGameOver) {
			squash.Reset(cfg)
			squash.State = app.StatePlaying
		}

		args[0].Call("preventDefault")
		return nil
	}))

	// Pause/Resume
	canvas.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) any {
		args[0].Call("preventDefault")

		if squash.State == app.StatePlaying {
			squash.State = app.StatePaused
		} else if squash.State == app.StatePaused {
			squash.State = app.StatePlaying
		}
		return nil
	}))

	canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) any {
		if squash.State != app.StatePlaying {
			return nil
		}

		rect := canvas.Call("getBoundingClientRect")
		mouseY := args[0].Get("clientY").Float() - rect.Get("top").Float()
		squash.CalcMovePaddle(mouseY)
		return nil
	}))
}
