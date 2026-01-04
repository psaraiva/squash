package main

import (
	"syscall/js"
	"time"

	"github.com/psaraiva/squash/internal/app"
	"github.com/psaraiva/squash/internal/ports"
	inputwasm "github.com/psaraiva/squash/pkg/adapters/input/wasm"
	inputweb "github.com/psaraiva/squash/pkg/adapters/input/web"
	outputweb "github.com/psaraiva/squash/pkg/adapters/output/web"
)

func main() {
	doc := js.Global().Get("document")
	canvasElement := doc.Call("getElementById", "gameCanvas")
	if canvasElement.IsNull() {
		panic("Canvas não encontrado")
	}

	w := canvasElement.Get("width").Float()
	h := canvasElement.Get("height").Float()

	var loader ports.ConfigProvider = inputwasm.NewConfigLoader()
	cfg := loader.Load()

	deltaTime := getDeltaTime(cfg.Fps)
	cfg.DeltaTime = float64(deltaTime) / 1000.0

	squash := app.NewSquash(w, h, cfg)
	inputwasm.SetupMouseHandlers(squash, canvasElement, cfg)

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(deltaTime) * time.Millisecond)
		defer ticker.Stop()

		var renderer ports.Renderer = outputweb.NewRenderer()
		for range ticker.C {
			squash.Update()
			inputweb.PaintGame(renderer, squash)
		}
	}()
	<-done
}

func getDeltaTime(fps int) int {
	deltaTime := 33 // 30 fps
	if fps == 60 {
		deltaTime = 16
	}

	return deltaTime
}
