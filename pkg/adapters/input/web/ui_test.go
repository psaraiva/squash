package web

import (
	"testing"

	"github.com/psaraiva/squash/internal/app"
	"github.com/psaraiva/squash/internal/ports/mocks"

	"github.com/stretchr/testify/mock"
)

func TestPaintGameStateMenu(t *testing.T) {
	tests := []struct {
		name              string
		gameState         app.GameState
		width             float64
		height            float64
		score             int
		lives             int
		textWidth         float64
		isTouchEnabled    bool
		expectBall        bool
		expectPaddle      bool
		minTextCalls      int
		expectMeasureText bool
	}{
		{
			name:              "Menu state - mouse controls",
			gameState:         app.StateMenu,
			width:             800.0,
			height:            600.0,
			score:             0,
			lives:             3,
			textWidth:         200.0,
			isTouchEnabled:    false,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      4, // Score, Lives, Menu text (2 lines)
			expectMeasureText: true,
		},
		{
			name:              "Menu state - touch controls",
			gameState:         app.StateMenu,
			width:             1024.0,
			height:            768.0,
			score:             50,
			lives:             2,
			textWidth:         250.0,
			isTouchEnabled:    true,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      4,
			expectMeasureText: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("Clear").Return()
			mockRenderer.On("DrawText", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("MeasureText", mock.Anything).Return(tt.textWidth)

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   tt.lives,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)
			g.State = tt.gameState
			g.Score = tt.score

			PaintGame(mockRenderer, g)

			mockRenderer.AssertCalled(t, "Clear")
			if tt.expectMeasureText {
				mockRenderer.AssertCalled(t, "MeasureText", mock.Anything)
			}
		})
	}
}

func TestPaintGameStatePlaying(t *testing.T) {
	tests := []struct {
		name            string
		gameState       app.GameState
		width           float64
		height          float64
		score           int
		lives           int
		expectBall      bool
		expectPaddle    bool
		expectDebugInfo bool
		debugMode       bool
	}{
		{
			name:            "Playing state - normal mode",
			gameState:       app.StatePlaying,
			width:           800.0,
			height:          600.0,
			score:           100,
			lives:           3,
			expectBall:      true,
			expectPaddle:    true,
			expectDebugInfo: false,
			debugMode:       false,
		},
		{
			name:            "Playing state - debug mode",
			gameState:       app.StatePlaying,
			width:           800.0,
			height:          600.0,
			score:           200,
			lives:           2,
			expectBall:      true,
			expectPaddle:    true,
			expectDebugInfo: true,
			debugMode:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("Clear").Return()
			mockRenderer.On("DrawText", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("DrawBall", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("DrawPaddle", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			if tt.debugMode {
				mockRenderer.On("DrawDebugText", mock.Anything, mock.Anything, mock.Anything).Return()
			}

			cfg := app.Config{
				Debug:          tt.debugMode,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   tt.lives,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)
			g.State = tt.gameState
			g.Score = tt.score

			PaintGame(mockRenderer, g)

			mockRenderer.AssertCalled(t, "Clear")
			if tt.expectBall {
				mockRenderer.AssertCalled(t, "DrawBall", mock.Anything, mock.Anything, mock.Anything)
			}
			if tt.expectPaddle {
				mockRenderer.AssertCalled(t, "DrawPaddle", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			if tt.expectDebugInfo {
				mockRenderer.AssertCalled(t, "DrawDebugText", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestPaintGameStatePaused(t *testing.T) {
	tests := []struct {
		name              string
		gameState         app.GameState
		width             float64
		height            float64
		textWidth         float64
		isTouchEnabled    bool
		expectBall        bool
		expectPaddle      bool
		minTextCalls      int
		expectMeasureText bool
	}{
		{
			name:              "Paused state - mouse controls",
			gameState:         app.StatePaused,
			width:             800.0,
			height:            600.0,
			textWidth:         200.0,
			isTouchEnabled:    false,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      3, // Score, Lives, Paused text
			expectMeasureText: true,
		},
		{
			name:              "Paused state - touch controls",
			gameState:         app.StatePaused,
			width:             1024.0,
			height:            768.0,
			textWidth:         250.0,
			isTouchEnabled:    true,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      3,
			expectMeasureText: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("Clear").Return()
			mockRenderer.On("DrawText", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("MeasureText", mock.Anything).Return(tt.textWidth)

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)
			g.State = tt.gameState

			PaintGame(mockRenderer, g)

			mockRenderer.AssertCalled(t, "Clear")
			if tt.expectMeasureText {
				mockRenderer.AssertCalled(t, "MeasureText", mock.Anything)
			}
		})
	}
}

func TestPaintGameStateGameOver(t *testing.T) {
	tests := []struct {
		name              string
		gameState         app.GameState
		width             float64
		height            float64
		score             int
		textWidth         float64
		isTouchEnabled    bool
		expectBall        bool
		expectPaddle      bool
		minTextCalls      int
		expectMeasureText bool
	}{
		{
			name:              "Game Over state - mouse controls",
			gameState:         app.StateGameOver,
			width:             800.0,
			height:            600.0,
			score:             500,
			textWidth:         200.0,
			isTouchEnabled:    false,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      4, // Score, Lives, Game Over text, Restart text
			expectMeasureText: true,
		},
		{
			name:              "Game Over state - touch controls",
			gameState:         app.StateGameOver,
			width:             1024.0,
			height:            768.0,
			score:             1000,
			textWidth:         250.0,
			isTouchEnabled:    true,
			expectBall:        false,
			expectPaddle:      false,
			minTextCalls:      4,
			expectMeasureText: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("Clear").Return()
			mockRenderer.On("DrawText", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("MeasureText", mock.Anything).Return(tt.textWidth)

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   0, // Game over state
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)
			g.State = tt.gameState
			g.Score = tt.score

			PaintGame(mockRenderer, g)

			mockRenderer.AssertCalled(t, "Clear")
			if tt.expectMeasureText {
				mockRenderer.AssertCalled(t, "MeasureText", mock.Anything)
			}
		})
	}
}

func TestDrawGameElements(t *testing.T) {
	tests := []struct {
		name     string
		ballX    float64
		ballY    float64
		ballSize float64
		paddleX  float64
		paddleY  float64
		paddleW  float64
		paddleH  float64
	}{
		{
			name:     "Draw game elements - standard position",
			ballX:    400.0,
			ballY:    300.0,
			ballSize: 15.0,
			paddleX:  10.0,
			paddleY:  270.0,
			paddleW:  10.0,
			paddleH:  60.0,
		},
		{
			name:     "Draw game elements - corner position",
			ballX:    100.0,
			ballY:    100.0,
			ballSize: 10.0,
			paddleX:  5.0,
			paddleY:  50.0,
			paddleW:  15.0,
			paddleH:  80.0,
		},
		{
			name:     "Draw game elements - different sizes",
			ballX:    500.0,
			ballY:    400.0,
			ballSize: 20.0,
			paddleX:  20.0,
			paddleY:  350.0,
			paddleW:  12.0,
			paddleH:  70.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("DrawBall", tt.ballX, tt.ballY, tt.ballSize).Return()
			mockRenderer.On("DrawPaddle", tt.paddleX, tt.paddleY, tt.paddleW, tt.paddleH).Return()

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)
			g.BallX = tt.ballX
			g.BallY = tt.ballY
			g.BallSize = tt.ballSize
			g.PaddleX = tt.paddleX
			g.PaddleY = tt.paddleY
			g.PaddleW = tt.paddleW
			g.PaddleH = tt.paddleH

			drawGameElements(mockRenderer, g)

			mockRenderer.AssertCalled(t, "DrawBall", tt.ballX, tt.ballY, tt.ballSize)
			mockRenderer.AssertCalled(t, "DrawPaddle", tt.paddleX, tt.paddleY, tt.paddleW, tt.paddleH)
		})
	}
}

func TestDrawDebugInfo(t *testing.T) {
	tests := []struct {
		name               string
		fps                int
		lastLevel          int
		ballSize           float64
		spawnX             float64
		spawnY             float64
		ballX              float64
		ballY              float64
		ballDX             float64
		ballDY             float64
		expectedDebugCalls int
	}{
		{
			name:               "Debug info - level 0",
			fps:                60,
			lastLevel:          0,
			ballSize:           15.0,
			spawnX:             400.0,
			spawnY:             300.0,
			ballX:              450.0,
			ballY:              320.0,
			ballDX:             300.0,
			ballDY:             -300.0,
			expectedDebugCalls: 8,
		},
		{
			name:               "Debug info - level 5",
			fps:                120,
			lastLevel:          5,
			ballSize:           20.0,
			spawnX:             500.0,
			spawnY:             400.0,
			ballX:              550.0,
			ballY:              420.0,
			ballDX:             -450.0,
			ballDY:             450.0,
			expectedDebugCalls: 8,
		},
		{
			name:               "Debug info - negative velocities",
			fps:                30,
			lastLevel:          10,
			ballSize:           10.0,
			spawnX:             200.0,
			spawnY:             150.0,
			ballX:              180.0,
			ballY:              140.0,
			ballDX:             -600.0,
			ballDY:             -600.0,
			expectedDebugCalls: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("DrawDebugText", mock.Anything, mock.Anything, mock.Anything).Return()

			cfg := app.Config{
				Debug:          true,
				Fps:            tt.fps,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   tt.lastLevel,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)
			g.BallSize = tt.ballSize
			g.BallSpawnX = tt.spawnX
			g.BallSpawnY = tt.spawnY
			g.BallX = tt.ballX
			g.BallY = tt.ballY
			g.BallDX = tt.ballDX
			g.BallDY = tt.ballDY

			debugInfo := getDebugInfo(g)
			drawDebugInfo(mockRenderer, debugInfo)

			mockRenderer.AssertNumberOfCalls(t, "DrawDebugText", tt.expectedDebugCalls)
		})
	}
}

func TestPaintGameTextPositioning(t *testing.T) {
	tests := []struct {
		name      string
		state     app.GameState
		width     float64
		height    float64
		textWidth float64
	}{
		{
			name:      "Text centering - menu state",
			state:     app.StateMenu,
			width:     800.0,
			height:    600.0,
			textWidth: 200.0,
		},
		{
			name:      "Text centering - game over state",
			state:     app.StateGameOver,
			width:     1024.0,
			height:    768.0,
			textWidth: 250.0,
		},
		{
			name:      "Text centering - paused state",
			state:     app.StatePaused,
			width:     640.0,
			height:    480.0,
			textWidth: 180.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("Clear").Return()
			mockRenderer.On("DrawText", mock.Anything, mock.Anything, mock.Anything).Return()
			mockRenderer.On("MeasureText", mock.Anything).Return(tt.textWidth)

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)
			g.State = tt.state

			PaintGame(mockRenderer, g)

			mockRenderer.AssertCalled(t, "Clear")
		})
	}
}

func TestGetTextScore(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{
			name:  "Score zero",
			score: 0,
			want:  "Score: 0",
		},
		{
			name:  "Score positive",
			score: 100,
			want:  "Score: 100",
		},
		{
			name:  "Score high value",
			score: 999999,
			want:  "Score: 999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)
			g.Score = tt.score

			got := getTextScore(g)
			if got != tt.want {
				t.Errorf("getTextScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTextLives(t *testing.T) {
	tests := []struct {
		name  string
		lives int
		want  string
	}{
		{
			name:  "Lives zero",
			lives: 0,
			want:  "Lives: 0",
		},
		{
			name:  "Lives standard",
			lives: 3,
			want:  "Lives: 3",
		},
		{
			name:  "Lives high value",
			lives: 99,
			want:  "Lives: 99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   tt.lives,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)

			got := getTextLives(g)
			if got != tt.want {
				t.Errorf("getTextLives() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDebugInfoLines(t *testing.T) {
	tests := []struct {
		name      string
		fps       int
		level     int
		ballSize  float64
		spawnX    float64
		spawnY    float64
		ballX     float64
		ballY     float64
		ballDX    float64
		ballDY    float64
		wantLines int
	}{
		{
			name:      "Debug info standard",
			fps:       60,
			level:     5,
			ballSize:  15.0,
			spawnX:    400.0,
			spawnY:    300.0,
			ballX:     450.0,
			ballY:     320.0,
			ballDX:    300.0,
			ballDY:    -300.0,
			wantLines: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := app.Config{
				Debug:          true,
				Fps:            tt.fps,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   tt.level,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)
			g.BallSize = tt.ballSize
			g.BallSpawnX = tt.spawnX
			g.BallSpawnY = tt.spawnY
			g.BallX = tt.ballX
			g.BallY = tt.ballY
			g.BallDX = tt.ballDX
			g.BallDY = tt.ballDY

			got := getDebugInfo(g)
			if len(got) != tt.wantLines {
				t.Errorf("getDebugInfo() lines = %v, want %v", len(got), tt.wantLines)
			}
		})
	}
}

func TestDrawTextScore(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Draw score text",
			text: "Score: 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("DrawText", tt.text, 30.0, 30.0).Return()

			drawTextScore(mockRenderer, tt.text)

			mockRenderer.AssertCalled(t, "DrawText", tt.text, 30.0, 30.0)
		})
	}
}

func TestDrawTextLives(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		width  float64
		height float64
	}{
		{
			name:   "Draw lives text",
			text:   "Lives: 3",
			width:  800.0,
			height: 600.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("DrawText", tt.text, tt.width-120, 30.0).Return()

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)

			drawTextLives(mockRenderer, g, tt.text)

			mockRenderer.AssertCalled(t, "DrawText", tt.text, tt.width-120, 30.0)
		})
	}
}

func TestDrawDebugInfoCall(t *testing.T) {
	tests := []struct {
		name      string
		debugInfo []string
	}{
		{
			name: "Draw debug info",
			debugInfo: []string{
				"Game:.....",
				"FPS:      60",
				"Level:    5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			for i := range tt.debugInfo {
				mockRenderer.On("DrawDebugText", mock.Anything, 10.0, 80+float64(i*15)).Return()
			}

			drawDebugInfo(mockRenderer, tt.debugInfo)

			for i := range tt.debugInfo {
				mockRenderer.AssertCalled(t, "DrawDebugText", mock.Anything, 10.0, 80+float64(i*15))
			}
		})
	}
}

func TestDrawTextCenter(t *testing.T) {
	tests := []struct {
		name      string
		text      []string
		width     float64
		height    float64
		textWidth float64
	}{
		{
			name:      "Draw centered text",
			text:      []string{"Line 1", "Line 2"},
			width:     800.0,
			height:    600.0,
			textWidth: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			for i := range tt.text {
				mockRenderer.On("MeasureText", tt.text[i]).Return(tt.textWidth)
				mockRenderer.On("DrawText", tt.text[i], (tt.width-tt.textWidth)/2, tt.height/2+float64(i*20)).Return()
			}

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(tt.width, tt.height, cfg)

			drawTextCenter(mockRenderer, g, tt.text)

			for i := range tt.text {
				mockRenderer.AssertCalled(t, "MeasureText", tt.text[i])
				mockRenderer.AssertCalled(t, "DrawText", tt.text[i], (tt.width-tt.textWidth)/2, tt.height/2+float64(i*20))
			}
		})
	}
}

func TestDrawGameElementsCall(t *testing.T) {
	tests := []struct {
		name    string
		ballX   float64
		ballY   float64
		ballS   float64
		paddleX float64
		paddleY float64
		paddleW float64
		paddleH float64
	}{
		{
			name:    "Draw game elements call",
			ballX:   400.0,
			ballY:   300.0,
			ballS:   15.0,
			paddleX: 10.0,
			paddleY: 270.0,
			paddleW: 10.0,
			paddleH: 60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRenderer := mocks.NewRenderer(t)
			mockRenderer.On("DrawBall", tt.ballX, tt.ballY, tt.ballS).Return()
			mockRenderer.On("DrawPaddle", tt.paddleX, tt.paddleY, tt.paddleW, tt.paddleH).Return()

			cfg := app.Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			g := app.NewSquash(800, 600, cfg)
			g.BallX = tt.ballX
			g.BallY = tt.ballY
			g.BallSize = tt.ballS
			g.PaddleX = tt.paddleX
			g.PaddleY = tt.paddleY
			g.PaddleW = tt.paddleW
			g.PaddleH = tt.paddleH

			drawGameElements(mockRenderer, g)

			mockRenderer.AssertCalled(t, "DrawBall", tt.ballX, tt.ballY, tt.ballS)
			mockRenderer.AssertCalled(t, "DrawPaddle", tt.paddleX, tt.paddleY, tt.paddleW, tt.paddleH)
		})
	}
}
