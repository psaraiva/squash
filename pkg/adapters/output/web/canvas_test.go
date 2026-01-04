//go:build js && wasm

package web

import (
	"syscall/js"
	"testing"

	"github.com/psaraiva/squash/internal/ports"
)

// mockCall represents an expected method call
type mockCall struct {
	method string
	args   []interface{}
}

// mockJSContext is a test-specific mock implementation of JSContext interface
type mockJSContext struct {
	setCalls      []mockCall
	callCalls     []mockCall
	getCalls      []mockCall
	floatReturn   float64
	returnContext JSContext
	t             *testing.T
}

func newMockJSContext(t *testing.T) *mockJSContext {
	return &mockJSContext{
		setCalls:      make([]mockCall, 0),
		callCalls:     make([]mockCall, 0),
		getCalls:      make([]mockCall, 0),
		returnContext: nil,
		t:             t,
	}
}

func (m *mockJSContext) expectSet(key string, value interface{}) *mockJSContext {
	m.setCalls = append(m.setCalls, mockCall{method: "Set", args: []interface{}{key, value}})
	return m
}

func (m *mockJSContext) expectCall(method string, args ...interface{}) *mockJSContext {
	m.callCalls = append(m.callCalls, mockCall{method: "Call", args: append([]interface{}{method}, args...)})
	return m
}

func (m *mockJSContext) expectGet(key string) *mockJSContext {
	m.getCalls = append(m.getCalls, mockCall{method: "Get", args: []interface{}{key}})
	return m
}

func (m *mockJSContext) withFloatReturn(value float64) *mockJSContext {
	m.floatReturn = value
	return m
}

func (m *mockJSContext) withReturnContext(ctx JSContext) *mockJSContext {
	m.returnContext = ctx
	return m
}

func (m *mockJSContext) Set(key string, value interface{}) {
	if len(m.setCalls) == 0 {
		m.t.Errorf("Unexpected Set call: Set(%v, %v)", key, value)
		return
	}
	expected := m.setCalls[0]
	m.setCalls = m.setCalls[1:]

	if expected.args[0] != key || expected.args[1] != value {
		m.t.Errorf("Set() args = (%v, %v), want (%v, %v)", key, value, expected.args[0], expected.args[1])
	}
}

func (m *mockJSContext) Call(method string, args ...interface{}) JSContext {
	if len(m.callCalls) == 0 {
		m.t.Errorf("Unexpected Call: Call(%v, %v)", method, args)
		return m
	}
	m.callCalls = m.callCalls[1:]

	if m.returnContext != nil {
		return m.returnContext
	}
	return m
}

func (m *mockJSContext) Get(key string) JSContext {
	if len(m.getCalls) == 0 {
		m.t.Errorf("Unexpected Get call: Get(%v)", key)
		return m
	}
	m.getCalls = m.getCalls[1:]

	if m.returnContext != nil {
		return m.returnContext
	}
	return m
}

func (m *mockJSContext) Float() float64 {
	return m.floatReturn
}

func (m *mockJSContext) verify() {
	if len(m.setCalls) > 0 {
		m.t.Errorf("Expected Set calls not made: %v", m.setCalls)
	}
	if len(m.callCalls) > 0 {
		m.t.Errorf("Expected Call calls not made: %v", m.callCalls)
	}
	if len(m.getCalls) > 0 {
		m.t.Errorf("Expected Get calls not made: %v", m.getCalls)
	}
}

func TestNewCanvas(t *testing.T) {
	tests := []struct {
		name   string
		width  float64
		height float64
	}{
		{
			name:   "Standard canvas size",
			width:  800.0,
			height: 600.0,
		},
		{
			name:   "Large canvas size",
			width:  1920.0,
			height: 1080.0,
		},
		{
			name:   "Small canvas size",
			width:  640.0,
			height: 480.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCanvas := js.ValueOf(map[string]interface{}{
				"width":  tt.width,
				"height": tt.height,
			})

			canvas := &Canvas{
				ctx: NewJSContext(mockCanvas),
				w:   tt.width,
				h:   tt.height,
			}

			if canvas.w != tt.width {
				t.Errorf("NewCanvas() width = %v, want %v", canvas.w, tt.width)
			}
			if canvas.h != tt.height {
				t.Errorf("NewCanvas() height = %v, want %v", canvas.h, tt.height)
			}
		})
	}
}

func TestCanvasClear(t *testing.T) {
	tests := []struct {
		name   string
		width  float64
		height float64
	}{
		{
			name:   "Clear standard canvas",
			width:  800.0,
			height: 600.0,
		},
		{
			name:   "Clear large canvas",
			width:  1920.0,
			height: 1080.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := newMockJSContext(t)
			mockCtx.expectSet("fillStyle", "black")
			mockCtx.expectCall("fillRect")

			canvas := &Canvas{
				ctx: mockCtx,
				w:   tt.width,
				h:   tt.height,
			}

			canvas.Clear()

			mockCtx.verify()
		})
	}
}

func TestCanvasDrawBall(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		size float64
	}{
		{
			name: "Ball at origin",
			x:    0.0,
			y:    0.0,
			size: 10.0,
		},
		{
			name: "Ball at center",
			x:    400.0,
			y:    300.0,
			size: 15.0,
		},
		{
			name: "Ball with large size",
			x:    100.0,
			y:    100.0,
			size: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := newMockJSContext(t)
			mockCtx.expectSet("fillStyle", "white")
			mockCtx.expectCall("fillRect")

			canvas := &Canvas{
				ctx: mockCtx,
				w:   800.0,
				h:   600.0,
			}

			canvas.DrawBall(tt.x, tt.y, tt.size)

			mockCtx.verify()
		})
	}
}

func TestCanvasDrawPaddle(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		w    float64
		h    float64
	}{
		{
			name: "Standard paddle",
			x:    10.0,
			y:    270.0,
			w:    10.0,
			h:    60.0,
		},
		{
			name: "Wide paddle",
			x:    5.0,
			y:    250.0,
			w:    15.0,
			h:    80.0,
		},
		{
			name: "Tall paddle",
			x:    10.0,
			y:    200.0,
			w:    8.0,
			h:    100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := newMockJSContext(t)
			mockCtx.expectSet("fillStyle", "white")
			mockCtx.expectCall("fillRect")

			canvas := &Canvas{
				ctx: mockCtx,
				w:   800.0,
				h:   600.0,
			}

			canvas.DrawPaddle(tt.x, tt.y, tt.w, tt.h)

			mockCtx.verify()
		})
	}
}

func TestCanvasDrawText(t *testing.T) {
	tests := []struct {
		name string
		text string
		x    float64
		y    float64
	}{
		{
			name: "Score text",
			text: "Score: 100",
			x:    30.0,
			y:    30.0,
		},
		{
			name: "Lives text",
			text: "Lives: 3",
			x:    680.0,
			y:    30.0,
		},
		{
			name: "Game over text",
			text: "GAME OVER",
			x:    300.0,
			y:    300.0,
		},
		{
			name: "Empty text",
			text: "",
			x:    100.0,
			y:    100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := newMockJSContext(t)
			mockCtx.expectSet("fillStyle", "white")
			mockCtx.expectSet("font", "20px Arial")
			mockCtx.expectCall("fillText")

			canvas := &Canvas{
				ctx: mockCtx,
				w:   800.0,
				h:   600.0,
			}

			canvas.DrawText(tt.text, tt.x, tt.y)

			mockCtx.verify()
		})
	}
}

func TestCanvasDrawDebugText(t *testing.T) {
	tests := []struct {
		name string
		text string
		x    float64
		y    float64
	}{
		{
			name: "FPS debug info",
			text: "FPS: 60",
			x:    10.0,
			y:    80.0,
		},
		{
			name: "Level debug info",
			text: "Level: 5",
			x:    10.0,
			y:    95.0,
		},
		{
			name: "Position debug info",
			text: "Position: [400.0, 300.0]",
			x:    10.0,
			y:    110.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := newMockJSContext(t)
			mockCtx.expectSet("font", "12px monospace")
			mockCtx.expectSet("fillStyle", "#FFFF00")
			mockCtx.expectCall("fillText")

			canvas := &Canvas{
				ctx: mockCtx,
				w:   800.0,
				h:   600.0,
			}

			canvas.DrawDebugText(tt.text, tt.x, tt.y)

			mockCtx.verify()
		})
	}
}

func TestCanvasMeasureText(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedWidth float64
	}{
		{
			name:          "Short text",
			text:          "Hi",
			expectedWidth: 50.0,
		},
		{
			name:          "Medium text",
			text:          "Score: 100",
			expectedWidth: 120.0,
		},
		{
			name:          "Long text",
			text:          "GAME OVER - Press SPACE",
			expectedWidth: 250.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMetrics := newMockJSContext(t)
			mockMetrics.expectGet("width")
			mockMetrics.withFloatReturn(tt.expectedWidth)

			mockCtx := newMockJSContext(t)
			mockCtx.expectCall("measureText")
			mockCtx.withReturnContext(mockMetrics)

			canvas := &Canvas{
				ctx: mockCtx,
				w:   800.0,
				h:   600.0,
			}

			width := canvas.MeasureText(tt.text)

			if width != tt.expectedWidth {
				t.Errorf("MeasureText() = %v, want %v", width, tt.expectedWidth)
			}

			mockCtx.verify()
			mockMetrics.verify()
		})
	}
}

func TestCanvasImplementsRenderer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Canvas implements Renderer interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var _ ports.Renderer = (*Canvas)(nil)
		})
	}
}
