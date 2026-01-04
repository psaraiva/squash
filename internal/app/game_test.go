package app

import (
	"testing"
)

func TestCalcSpeedFactor(t *testing.T) {
	tests := []struct {
		name      string
		level     int
		increment float64
		want      float64
	}{
		{
			name:      "Level 0 - Base speed",
			level:     0,
			increment: 0.5,
			want:      1.0,
		},
		{
			name:      "Level 1 - First increment",
			level:     1,
			increment: 0.5,
			want:      1.5,
		},
		{
			name:      "Level 5 - Mid level",
			level:     5,
			increment: 0.5,
			want:      3.5,
		},
		{
			name:      "Level 10 - High level",
			level:     10,
			increment: 0.5,
			want:      6.0,
		},
		{
			name:      "Level negative - Should clamp to 0",
			level:     -1,
			increment: 0.5,
			want:      1.0,
		},
		{
			name:      "Level 100 - Should clamp to 99",
			level:     100,
			increment: 0.5,
			want:      1.0,
		},
		{
			name:      "Increment negative - Should return 1.0",
			level:     5,
			increment: -0.1,
			want:      1.0,
		},
		{
			name:      "Increment above 1.0 - Should return 1.0",
			level:     5,
			increment: 1.5,
			want:      1.0,
		},
		{
			name:      "Zero increment - No speed increase",
			level:     10,
			increment: 0.0,
			want:      1.0,
		},
		{
			name:      "Max valid increment (1.0)",
			level:     5,
			increment: 1.0,
			want:      6.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcSpeedFactor(tt.level, tt.increment)
			if got != tt.want {
				t.Errorf("calcSpeedFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcBallSize(t *testing.T) {
	tests := []struct {
		name  string
		scale float64
		want  float64
	}{
		{
			name:  "No scaling (0.0)",
			scale: 0.0,
			want:  10.0,
		},
		{
			name:  "Half scale (0.5)",
			scale: 0.5,
			want:  15.0,
		},
		{
			name:  "Full scale (1.0)",
			scale: 1.0,
			want:  20.0,
		},
		{
			name:  "Quarter scale (0.25)",
			scale: 0.25,
			want:  12.5,
		},
		{
			name:  "Negative scale - Should return base",
			scale: -0.5,
			want:  10.0,
		},
		{
			name:  "Above max scale - Should return base",
			scale: 1.5,
			want:  10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcBallSize(tt.scale)
			if got != tt.want {
				t.Errorf("calcBallSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcBallStartY(t *testing.T) {
	tests := []struct {
		name         string
		height       float64
		randomSource float64
		want         float64
	}{
		{
			name:         "Lower boundary (Random 0.0)",
			height:       600.0,
			randomSource: 0.0,
			want:         90.0, // 15% of 600
		},
		{
			name:         "Upper boundary (Random 1.0)",
			height:       600.0,
			randomSource: 1.0,
			want:         510.0, // 90 (margin) + 420 (70% of 600)
		},
		{
			name:         "Exact center (Random 0.5)",
			height:       600.0,
			randomSource: 0.5,
			want:         300.0,
		},
		{
			name:         "Smaller board (Height 100)",
			height:       100.0,
			randomSource: 0.5,
			want:         50.0, // Margin 15 + (0.5 * 70) = 50
		},
		{
			name:         "Large board (Height 1000)",
			height:       1000.0,
			randomSource: 0.25,
			want:         325.0, // 150 (margin) + (0.25 * 700)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcBallStartY(tt.height, tt.randomSource)
			if got != tt.want {
				t.Errorf("calcBallStartY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcRespawnPaddleY(t *testing.T) {
	tests := []struct {
		name    string
		height  float64
		paddleH float64
		want    float64
	}{
		{
			name:    "Standard board - 600x60 paddle",
			height:  600.0,
			paddleH: 60.0,
			want:    270.0, // (600 / 2) - (60 / 2)
		},
		{
			name:    "Small board - 100x20 paddle",
			height:  100.0,
			paddleH: 20.0,
			want:    40.0,
		},
		{
			name:    "Large board - 1000x80 paddle",
			height:  1000.0,
			paddleH: 80.0,
			want:    460.0,
		},
		{
			name:    "Zero height paddle",
			height:  600.0,
			paddleH: 0.0,
			want:    300.0,
		},
		{
			name:    "Paddle equals height",
			height:  100.0,
			paddleH: 100.0,
			want:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcRespawPaddleY(tt.height, tt.paddleH)
			if got != tt.want {
				t.Errorf("calcRespawnPaddleY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcRandomDirectionStartBall(t *testing.T) {
	tests := []struct {
		name   string
		ballDX float64
		ballDY float64
	}{
		{
			name:   "Positive values",
			ballDX: 300.0,
			ballDY: 300.0,
		},
		{
			name:   "Negative values",
			ballDX: -300.0,
			ballDY: -300.0,
		},
		{
			name:   "Zero values",
			ballDX: 0.0,
			ballDY: 0.0,
		},
		{
			name:   "Mixed values - DX positive, DY negative",
			ballDX: 450.0,
			ballDY: -450.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDX, gotDY := calcRandomDirectionStartBall(tt.ballDX, tt.ballDY)

			// Check that the absolute values are preserved
			if absFloat(gotDX) != absFloat(tt.ballDX) {
				t.Errorf("calcRandomDirectionStartBall() DX magnitude changed: got %v, want %v", absFloat(gotDX), absFloat(tt.ballDX))
			}
			if absFloat(gotDY) != absFloat(tt.ballDY) {
				t.Errorf("calcRandomDirectionStartBall() DY magnitude changed: got %v, want %v", absFloat(gotDY), absFloat(tt.ballDY))
			}
		})
	}
}

func TestNewSquash(t *testing.T) {
	tests := []struct {
		name   string
		width  float64
		height float64
		cfg    Config
	}{
		{
			name:   "Standard initialization",
			width:  800.0,
			height: 600.0,
			cfg: Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			},
		},
		{
			name:   "Custom configuration",
			width:  1024.0,
			height: 768.0,
			cfg: Config{
				Debug:          true,
				Fps:            120,
				DeltaTime:      0.008,
				InitialLives:   5,
				InitialLevel:   5,
				SpeedIncrement: 0.3,
				BallScale:      1.0,
			},
		},
		{
			name:   "Minimal configuration",
			width:  640.0,
			height: 480.0,
			cfg: Config{
				Debug:          false,
				Fps:            30,
				DeltaTime:      0.033,
				InitialLives:   1,
				InitialLevel:   0,
				SpeedIncrement: 0.0,
				BallScale:      0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewSquash(tt.width, tt.height, tt.cfg)

			if game.Width != tt.width {
				t.Errorf("NewSquash() Width = %v, want %v", game.Width, tt.width)
			}
			if game.Height != tt.height {
				t.Errorf("NewSquash() Height = %v, want %v", game.Height, tt.height)
			}
			if game.State != StateMenu {
				t.Errorf("NewSquash() State = %v, want %v", game.State, StateMenu)
			}
			if game.PaddleH != 60 {
				t.Errorf("NewSquash() PaddleH = %v, want 60", game.PaddleH)
			}
			if game.PaddleW != 10 {
				t.Errorf("NewSquash() PaddleW = %v, want 10", game.PaddleW)
			}
			if game.PaddleX != 10 {
				t.Errorf("NewSquash() PaddleX = %v, want 10", game.PaddleX)
			}
			if game.DebugMode != tt.cfg.Debug {
				t.Errorf("NewSquash() DebugMode = %v, want %v", game.DebugMode, tt.cfg.Debug)
			}
			if game.Lives != tt.cfg.InitialLives {
				t.Errorf("NewSquash() Lives = %v, want %v", game.Lives, tt.cfg.InitialLives)
			}
		})
	}
}

func TestSquashReset(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "Reset with default config",
			cfg: Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			},
		},
		{
			name: "Reset with debug mode",
			cfg: Config{
				Debug:          true,
				Fps:            120,
				DeltaTime:      0.008,
				InitialLives:   5,
				InitialLevel:   10,
				SpeedIncrement: 0.8,
				BallScale:      1.0,
			},
		},
		{
			name: "Reset with high level",
			cfg: Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   2,
				InitialLevel:   20,
				SpeedIncrement: 0.4,
				BallScale:      0.75,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Squash{Width: 800, Height: 600, PaddleH: 60, PaddleW: 10}
			game.Reset(tt.cfg)

			if game.DebugMode != tt.cfg.Debug {
				t.Errorf("Reset() DebugMode = %v, want %v", game.DebugMode, tt.cfg.Debug)
			}
			if game.Fps != tt.cfg.Fps {
				t.Errorf("Reset() Fps = %v, want %v", game.Fps, tt.cfg.Fps)
			}
			if game.DeltaTime != tt.cfg.DeltaTime {
				t.Errorf("Reset() DeltaTime = %v, want %v", game.DeltaTime, tt.cfg.DeltaTime)
			}
			if game.Lives != tt.cfg.InitialLives {
				t.Errorf("Reset() Lives = %v, want %v", game.Lives, tt.cfg.InitialLives)
			}
			if game.LastLevel != tt.cfg.InitialLevel {
				t.Errorf("Reset() LastLevel = %v, want %v", game.LastLevel, tt.cfg.InitialLevel)
			}
			if game.SpeedIncrement != tt.cfg.SpeedIncrement {
				t.Errorf("Reset() SpeedIncrement = %v, want %v", game.SpeedIncrement, tt.cfg.SpeedIncrement)
			}
			if game.Score != tt.cfg.InitialLevel*100 {
				t.Errorf("Reset() Score = %v, want %v", game.Score, tt.cfg.InitialLevel*100)
			}
			expectedBallSize := calcBallSize(tt.cfg.BallScale)
			if game.BallSize != expectedBallSize {
				t.Errorf("Reset() BallSize = %v, want %v", game.BallSize, expectedBallSize)
			}
			expectedPaddleY := calcRespawPaddleY(game.Height, game.PaddleH)
			if game.PaddleY != expectedPaddleY {
				t.Errorf("Reset() PaddleY = %v, want %v", game.PaddleY, expectedPaddleY)
			}
		})
	}
}

func TestSquashRespawnBall(t *testing.T) {
	tests := []struct {
		name      string
		width     float64
		height    float64
		lastLevel int
		speedInc  float64
	}{
		{
			name:      "Basic respawn level 0",
			width:     800.0,
			height:    600.0,
			lastLevel: 0,
			speedInc:  0.5,
		},
		{
			name:      "High level respawn",
			width:     1024.0,
			height:    768.0,
			lastLevel: 10,
			speedInc:  0.8,
		},
		{
			name:      "Mid level respawn",
			width:     640.0,
			height:    480.0,
			lastLevel: 5,
			speedInc:  0.3,
		},
		{
			name:      "Zero increment",
			width:     800.0,
			height:    600.0,
			lastLevel: 3,
			speedInc:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Squash{
				Width:          tt.width,
				Height:         tt.height,
				LastLevel:      tt.lastLevel,
				SpeedIncrement: tt.speedInc,
			}

			game.respawnBall()

			if game.BallX != tt.width/2 {
				t.Errorf("respawnBall() BallX = %v, want %v", game.BallX, tt.width/2)
			}
			if game.BallX != game.BallSpawnX {
				t.Errorf("respawnBall() BallSpawnX = %v, want %v", game.BallSpawnX, game.BallX)
			}
			if game.BallY != game.BallSpawnY {
				t.Errorf("respawnBall() BallSpawnY = %v, want %v", game.BallSpawnY, game.BallY)
			}

			// Check ball is within playable area
			minY := tt.height * 0.15
			maxY := tt.height * 0.85
			if game.BallY < minY || game.BallY > maxY {
				t.Errorf("respawnBall() BallY = %v, out of range [%v, %v]", game.BallY, minY, maxY)
			}

			// Check speed magnitudes are set correctly
			expectedFactor := calcSpeedFactor(tt.lastLevel, tt.speedInc)
			expectedSpeed := BaseSpeedBall * expectedFactor
			if absFloat(game.BallDX) != expectedSpeed {
				t.Errorf("respawnBall() BallDX magnitude = %v, want %v", absFloat(game.BallDX), expectedSpeed)
			}
			if absFloat(game.BallDY) != expectedSpeed {
				t.Errorf("respawnBall() BallDY magnitude = %v, want %v", absFloat(game.BallDY), expectedSpeed)
			}
		})
	}
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestNewDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Default config values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewDefaultConfig()

			if cfg.Debug != false {
				t.Errorf("NewDefaultConfig() Debug = %v, want %v", cfg.Debug, false)
			}
			if cfg.InitialLives != 3 {
				t.Errorf("NewDefaultConfig() InitialLives = %v, want %v", cfg.InitialLives, 3)
			}
			if cfg.InitialLevel != 0 {
				t.Errorf("NewDefaultConfig() InitialLevel = %v, want %v", cfg.InitialLevel, 0)
			}
			if cfg.SpeedIncrement != 0.25 {
				t.Errorf("NewDefaultConfig() SpeedIncrement = %v, want %v", cfg.SpeedIncrement, 0.25)
			}
			if cfg.BallScale != 0.0 {
				t.Errorf("NewDefaultConfig() BallScale = %v, want %v", cfg.BallScale, 0.0)
			}
			if cfg.Fps != 30 {
				t.Errorf("NewDefaultConfig() Fps = %v, want %v", cfg.Fps, 30)
			}
		})
	}
}
