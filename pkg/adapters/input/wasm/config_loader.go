//go:build js && wasm

package wasm

import (
	"strconv"
	"syscall/js"

	"github.com/psaraiva/squash/internal/app"
	"github.com/psaraiva/squash/internal/ports"
)

type ConfigLoader struct {
	data []byte
}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (c *ConfigLoader) Load() app.Config {
	cfg := app.NewDefaultConfig()

	window := js.Global().Get("window")
	if window.IsUndefined() || window.IsNull() {
		return cfg
	}

	search := window.Get("location").Get("search")
	params := js.Global().Get("URLSearchParams").New(search)

	// 1. Debug
	if params.Call("has", "debug").Bool() {
		cfg.Debug = params.Call("get", "debug").String() == "true"
	}

	// 2. Lives (1 to 99)
	if params.Call("has", "lives").Bool() {
		cfg.InitialLives = 3
		if val, err := strconv.Atoi(params.Call("get", "lives").String()); err == nil {
			if val >= 1 && val <= 99 {
				cfg.InitialLives = val
			}
		}
	}

	// 3. Initial Level (0 to 50)
	if params.Call("has", "level").Bool() {
		cfg.InitialLevel = 0
		if val, err := strconv.Atoi(params.Call("get", "level").String()); err == nil {
			if val >= 0 && val <= 50 {
				cfg.InitialLevel = val
			}
		}
	}

	// 4. Speed Increment (0.0 to 1.0)
	if params.Call("has", "boost").Bool() {
		cfg.SpeedIncrement = 0.0
		if val, err := strconv.ParseFloat(params.Call("get", "boost").String(), 64); err == nil {
			if val >= 0.0 && val <= 1.0 {
				cfg.SpeedIncrement = val
			}
		}
	}

	// 5. Ball Size (0.0 to 1.0)
	if params.Call("has", "ballsize").Bool() {
		cfg.BallScale = 0.0
		if val, err := strconv.ParseFloat(params.Call("get", "ballsize").String(), 64); err == nil {
			if val >= 0.0 && val <= 1.0 {
				cfg.BallScale = val
			}
		}
	}

	// 6. FPS 30 or 60
	if params.Call("has", "fps").Bool() {
		cfg.Fps = 30
		if val, err := strconv.Atoi(params.Call("get", "fps").String()); err == nil {
			if val == 60 {
				cfg.Fps = val
			}
		}
	}

	return cfg
}

var _ ports.ConfigProvider = (*ConfigLoader)(nil)
