package app

import (
	"testing"
)

func TestSquashUpdate(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
	}{
		{
			name:  "Update when playing",
			state: StatePlaying,
		},
		{
			name:  "Update when in menu - should not update",
			state: StateMenu,
		},
		{
			name:  "Update when paused - should not update",
			state: StatePaused,
		},
		{
			name:  "Update when game over - should not update",
			state: StateGameOver,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.State = tt.state

			initialBallX := game.BallX
			initialBallY := game.BallY

			game.Update()

			if tt.state == StatePlaying {
				// When playing, ball should move
				if game.BallX == initialBallX && game.BallY == initialBallY {
					t.Errorf("Update() ball should have moved when state is %v", tt.state)
				}
			} else {
				// When not playing, ball should not move
				if game.BallX != initialBallX || game.BallY != initialBallY {
					t.Errorf("Update() ball should not have moved when state is %v", tt.state)
				}
			}
		})
	}
}

func TestCalcMovePaddle(t *testing.T) {
	tests := []struct {
		name        string
		axisY       float64
		state       GameState
		wantPaddleY float64
	}{
		{
			name:        "Move paddle to center",
			axisY:       300.0,
			state:       StatePlaying,
			wantPaddleY: 270.0, // 300 - (60/2)
		},
		{
			name:        "Move paddle to top boundary",
			axisY:       0.0,
			state:       StatePlaying,
			wantPaddleY: 0.0,
		},
		{
			name:        "Move paddle beyond top - should clamp",
			axisY:       -50.0,
			state:       StatePlaying,
			wantPaddleY: 0.0,
		},
		{
			name:        "Move paddle to bottom boundary",
			axisY:       600.0,
			state:       StatePlaying,
			wantPaddleY: 540.0, // 600 - 60
		},
		{
			name:        "Move paddle beyond bottom - should clamp",
			axisY:       650.0,
			state:       StatePlaying,
			wantPaddleY: 540.0,
		},
		{
			name:        "Move paddle when not playing - should not move",
			axisY:       300.0,
			state:       StateMenu,
			wantPaddleY: 270.0, // Original position (centered)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.State = tt.state

			game.CalcMovePaddle(tt.axisY)

			if tt.state == StatePlaying {
				if game.PaddleY != tt.wantPaddleY {
					t.Errorf("CalcMovePaddle() PaddleY = %v, want %v", game.PaddleY, tt.wantPaddleY)
				}
			} else {
				// When not playing, paddle should remain in original position
				expectedY := calcRespawPaddleY(game.Height, game.PaddleH)
				if game.PaddleY != expectedY {
					t.Errorf("CalcMovePaddle() PaddleY changed when not playing = %v, want %v", game.PaddleY, expectedY)
				}
			}
		})
	}
}

func TestSetPaddlePosition(t *testing.T) {
	tests := []struct {
		name        string
		posY        float64
		wantPaddleY float64
	}{
		{
			name:        "Set paddle to center",
			posY:        270.0,
			wantPaddleY: 270.0,
		},
		{
			name:        "Set paddle to top",
			posY:        0.0,
			wantPaddleY: 0.0,
		},
		{
			name:        "Set paddle beyond top - should clamp",
			posY:        -50.0,
			wantPaddleY: 0.0,
		},
		{
			name:        "Set paddle to bottom",
			posY:        540.0,
			wantPaddleY: 540.0,
		},
		{
			name:        "Set paddle beyond bottom - should clamp",
			posY:        600.0,
			wantPaddleY: 540.0,
		},
		{
			name:        "Set paddle to arbitrary valid position",
			posY:        150.0,
			wantPaddleY: 150.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)

			game.SetPaddlePosition(tt.posY)

			if game.PaddleY != tt.wantPaddleY {
				t.Errorf("SetPaddlePosition() PaddleY = %v, want %v", game.PaddleY, tt.wantPaddleY)
			}
		})
	}
}

func TestCalcCollisionWithWalls_TopAndBottom(t *testing.T) {
	tests := []struct {
		name          string
		ballY         float64
		ballDY        float64
		wantBallDY    float64
		shouldReverse bool
	}{
		{
			name:          "Collision at top",
			ballY:         0.0,
			ballDY:        -150.0,
			wantBallDY:    150.0,
			shouldReverse: true,
		},
		{
			name:          "Collision at bottom",
			ballY:         585.0, // 600 - 15 (ballSize)
			ballDY:        150.0,
			wantBallDY:    -150.0,
			shouldReverse: true,
		},
		{
			name:          "No collision - middle",
			ballY:         300.0,
			ballDY:        150.0,
			wantBallDY:    150.0,
			shouldReverse: false,
		},
		{
			name:          "No collision - near top but not touching",
			ballY:         10.0,
			ballDY:        150.0,
			wantBallDY:    150.0,
			shouldReverse: false,
		},
		{
			name:          "No collision - near bottom but not touching",
			ballY:         580.0,
			ballDY:        -150.0,
			wantBallDY:    -150.0,
			shouldReverse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.BallY = tt.ballY
			game.BallDY = tt.ballDY

			game.calcCollisionWithWalls()

			if game.BallDY != tt.wantBallDY {
				t.Errorf("calcCollisionWithWalls() BallDY = %v, want %v", game.BallDY, tt.wantBallDY)
			}
		})
	}
}

func TestCalcCollisionWithWalls_Right(t *testing.T) {
	tests := []struct {
		name          string
		ballX         float64
		ballDX        float64
		wantBallDX    float64
		wantBallX     float64
		shouldReverse bool
	}{
		{
			name:          "Collision at right wall",
			ballX:         785.0, // 800 - 15 (ballSize)
			ballDX:        300.0,
			wantBallDX:    -300.0,
			wantBallX:     784.0, // 800 - 15 - 1
			shouldReverse: true,
		},
		{
			name:          "Collision beyond right wall",
			ballX:         790.0,
			ballDX:        300.0,
			wantBallDX:    -300.0,
			wantBallX:     784.0,
			shouldReverse: true,
		},
		{
			name:          "No collision - middle",
			ballX:         400.0,
			ballDX:        300.0,
			wantBallDX:    300.0,
			wantBallX:     400.0,
			shouldReverse: false,
		},
		{
			name:          "No collision - near right but not touching",
			ballX:         780.0,
			ballDX:        300.0,
			wantBallDX:    300.0,
			wantBallX:     780.0,
			shouldReverse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.BallX = tt.ballX
			game.BallDX = tt.ballDX

			game.calcCollisionWithWalls()

			if game.BallDX != tt.wantBallDX {
				t.Errorf("calcCollisionWithWalls() BallDX = %v, want %v", game.BallDX, tt.wantBallDX)
			}
			if game.BallX != tt.wantBallX {
				t.Errorf("calcCollisionWithWalls() BallX = %v, want %v", game.BallX, tt.wantBallX)
			}
		})
	}
}

func TestCalcCollisionPaddle(t *testing.T) {
	tests := []struct {
		name          string
		ballX         float64
		ballY         float64
		ballDX        float64
		paddleY       float64
		initialScore  int
		wantBallDX    float64
		wantScore     int
		shouldCollide bool
	}{
		{
			name:          "Direct hit on paddle",
			ballX:         5.0,
			ballY:         270.0,
			ballDX:        -300.0,
			paddleY:       270.0,
			initialScore:  0,
			wantBallDX:    300.0,
			wantScore:     10,
			shouldCollide: true,
		},
		{
			name:          "Hit paddle at top edge",
			ballX:         5.0,
			ballY:         270.0,
			ballDX:        -300.0,
			paddleY:       260.0,
			initialScore:  50,
			wantBallDX:    300.0,
			wantScore:     60,
			shouldCollide: true,
		},
		{
			name:          "Hit paddle at bottom edge",
			ballX:         5.0,
			ballY:         315.0,
			ballDX:        -300.0,
			paddleY:       270.0,
			initialScore:  100,
			wantBallDX:    300.0,
			wantScore:     110,
			shouldCollide: true,
		},
		{
			name:          "Miss paddle - too high",
			ballX:         5.0,
			ballY:         200.0,
			ballDX:        -300.0,
			paddleY:       270.0,
			initialScore:  0,
			wantBallDX:    -300.0,
			wantScore:     0,
			shouldCollide: false,
		},
		{
			name:          "Miss paddle - too low",
			ballX:         5.0,
			ballY:         350.0,
			ballDX:        -300.0,
			paddleY:       270.0,
			initialScore:  0,
			wantBallDX:    -300.0,
			wantScore:     0,
			shouldCollide: false,
		},
		{
			name:          "Ball not at paddle X position",
			ballX:         100.0,
			ballY:         270.0,
			ballDX:        -300.0,
			paddleY:       270.0,
			initialScore:  0,
			wantBallDX:    -300.0,
			wantScore:     0,
			shouldCollide: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.BallX = tt.ballX
			game.BallY = tt.ballY
			game.BallDX = tt.ballDX
			game.PaddleY = tt.paddleY
			game.Score = tt.initialScore

			game.calcCollisionPaddle()

			if game.BallDX != tt.wantBallDX {
				t.Errorf("calcCollisionPaddle() BallDX = %v, want %v", game.BallDX, tt.wantBallDX)
			}
			if game.Score != tt.wantScore {
				t.Errorf("calcCollisionPaddle() Score = %v, want %v", game.Score, tt.wantScore)
			}
			if tt.shouldCollide {
				expectedBallX := game.PaddleX + game.PaddleW + 1
				if game.BallX != expectedBallX {
					t.Errorf("calcCollisionPaddle() BallX = %v, want %v", game.BallX, expectedBallX)
				}
			}
		})
	}
}

func TestCalcLostLive(t *testing.T) {
	tests := []struct {
		name          string
		ballX         float64
		initialLives  int
		wantLives     int
		wantState     GameState
		shouldRespawn bool
	}{
		{
			name:          "Lost live with lives remaining",
			ballX:         -20.0,
			initialLives:  3,
			wantLives:     2,
			wantState:     StatePlaying,
			shouldRespawn: true,
		},
		{
			name:          "Lost last live - game over",
			ballX:         -20.0,
			initialLives:  1,
			wantLives:     0,
			wantState:     StateGameOver,
			shouldRespawn: false,
		},
		{
			name:          "Lost live with 2 lives remaining",
			ballX:         -15.0,
			initialLives:  2,
			wantLives:     1,
			wantState:     StatePlaying,
			shouldRespawn: true,
		},
		{
			name:          "Ball not lost - positive X",
			ballX:         10.0,
			initialLives:  3,
			wantLives:     3,
			wantState:     StatePlaying,
			shouldRespawn: false,
		},
		{
			name:          "Ball at edge - not lost",
			ballX:         0.0,
			initialLives:  3,
			wantLives:     3,
			wantState:     StatePlaying,
			shouldRespawn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   tt.initialLives,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.State = StatePlaying
			game.BallX = tt.ballX
			game.Lives = tt.initialLives

			oldBallX := game.BallX
			game.calcLostLive()

			if game.Lives != tt.wantLives {
				t.Errorf("calcLostLive() Lives = %v, want %v", game.Lives, tt.wantLives)
			}
			if game.State != tt.wantState {
				t.Errorf("calcLostLive() State = %v, want %v", game.State, tt.wantState)
			}
			if tt.shouldRespawn {
				// Ball should be respawned (X position changed)
				if game.BallX == oldBallX {
					t.Errorf("calcLostLive() ball should have respawned")
				}
			}
		})
	}
}

func TestCalcNextLevel(t *testing.T) {
	tests := []struct {
		name           string
		score          int
		lastLevel      int
		ballDX         float64
		ballDY         float64
		speedIncrement float64
		wantLastLevel  int
		expectSpeedUp  bool
	}{
		{
			name:           "Reach level 1 from level 0",
			score:          100,
			lastLevel:      0,
			ballDX:         300.0,
			ballDY:         300.0,
			speedIncrement: 0.5,
			wantLastLevel:  1,
			expectSpeedUp:  true,
		},
		{
			name:           "Reach level 2 from level 1",
			score:          200,
			lastLevel:      1,
			ballDX:         450.0,
			ballDY:         -450.0,
			speedIncrement: 0.5,
			wantLastLevel:  2,
			expectSpeedUp:  true,
		},
		{
			name:           "Stay at same level - not enough score",
			score:          150,
			lastLevel:      1,
			ballDX:         450.0,
			ballDY:         450.0,
			speedIncrement: 0.5,
			wantLastLevel:  1,
			expectSpeedUp:  false,
		},
		{
			name:           "Reach level 5 from level 4",
			score:          500,
			lastLevel:      4,
			ballDX:         -750.0,
			ballDY:         750.0,
			speedIncrement: 0.5,
			wantLastLevel:  5,
			expectSpeedUp:  true,
		},
		{
			name:           "Level up with zero increment",
			score:          100,
			lastLevel:      0,
			ballDX:         300.0,
			ballDY:         300.0,
			speedIncrement: 0.0,
			wantLastLevel:  1,
			expectSpeedUp:  false,
		},
		{
			name:           "Multiple levels jump",
			score:          500,
			lastLevel:      0,
			ballDX:         300.0,
			ballDY:         -300.0,
			speedIncrement: 0.5,
			wantLastLevel:  5,
			expectSpeedUp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      0.016,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: tt.speedIncrement,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.Score = tt.score
			game.LastLevel = tt.lastLevel
			game.BallDX = tt.ballDX
			game.BallDY = tt.ballDY

			oldBallDX := game.BallDX
			oldBallDY := game.BallDY

			game.calcNextLevel()

			if game.LastLevel != tt.wantLastLevel {
				t.Errorf("calcNextLevel() LastLevel = %v, want %v", game.LastLevel, tt.wantLastLevel)
			}

			if tt.expectSpeedUp && tt.speedIncrement > 0 {
				// Speed should increase
				if absFloat(game.BallDX) <= absFloat(oldBallDX) {
					t.Errorf("calcNextLevel() BallDX speed should have increased: old=%v, new=%v", absFloat(oldBallDX), absFloat(game.BallDX))
				}
				if absFloat(game.BallDY) <= absFloat(oldBallDY) {
					t.Errorf("calcNextLevel() BallDY speed should have increased: old=%v, new=%v", absFloat(oldBallDY), absFloat(game.BallDY))
				}

				// Direction should be preserved
				if (oldBallDX > 0) != (game.BallDX > 0) {
					t.Errorf("calcNextLevel() BallDX direction changed")
				}
				if (oldBallDY > 0) != (game.BallDY > 0) {
					t.Errorf("calcNextLevel() BallDY direction changed")
				}
			} else if !tt.expectSpeedUp {
				// Speed should not change
				if game.BallDX != oldBallDX {
					t.Errorf("calcNextLevel() BallDX should not change: old=%v, new=%v", oldBallDX, game.BallDX)
				}
				if game.BallDY != oldBallDY {
					t.Errorf("calcNextLevel() BallDY should not change: old=%v, new=%v", oldBallDY, game.BallDY)
				}
			}
		})
	}
}

func TestCalcMoveBall(t *testing.T) {
	tests := []struct {
		name      string
		ballX     float64
		ballY     float64
		ballDX    float64
		ballDY    float64
		deltaTime float64
		wantBallX float64
		wantBallY float64
	}{
		{
			name:      "Move ball with positive velocities",
			ballX:     100.0,
			ballY:     200.0,
			ballDX:    300.0,
			ballDY:    300.0,
			deltaTime: 0.016,
			wantBallX: 104.8, // 100 + (300 * 0.016)
			wantBallY: 204.8, // 200 + (300 * 0.016)
		},
		{
			name:      "Move ball with negative velocities",
			ballX:     100.0,
			ballY:     200.0,
			ballDX:    -300.0,
			ballDY:    -300.0,
			deltaTime: 0.016,
			wantBallX: 95.2,  // 100 + (-300 * 0.016)
			wantBallY: 195.2, // 200 + (-300 * 0.016)
		},
		{
			name:      "Move ball with mixed velocities",
			ballX:     100.0,
			ballY:     200.0,
			ballDX:    300.0,
			ballDY:    -300.0,
			deltaTime: 0.016,
			wantBallX: 104.8,
			wantBallY: 195.2,
		},
		{
			name:      "Move ball with zero velocity",
			ballX:     100.0,
			ballY:     200.0,
			ballDX:    0.0,
			ballDY:    0.0,
			deltaTime: 0.016,
			wantBallX: 100.0,
			wantBallY: 200.0,
		},
		{
			name:      "Move ball with different delta time",
			ballX:     100.0,
			ballY:     200.0,
			ballDX:    300.0,
			ballDY:    300.0,
			deltaTime: 0.033,
			wantBallX: 109.9, // 100 + (300 * 0.033)
			wantBallY: 209.9, // 200 + (300 * 0.033)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Debug:          false,
				Fps:            60,
				DeltaTime:      tt.deltaTime,
				InitialLives:   3,
				InitialLevel:   0,
				SpeedIncrement: 0.5,
				BallScale:      0.5,
			}
			game := NewSquash(800, 600, cfg)
			game.BallX = tt.ballX
			game.BallY = tt.ballY
			game.BallDX = tt.ballDX
			game.BallDY = tt.ballDY

			game.calcMoveBall()

			// Use small epsilon for floating point comparison
			epsilon := 0.01
			if absFloat(game.BallX-tt.wantBallX) > epsilon {
				t.Errorf("calcMoveBall() BallX = %v, want %v", game.BallX, tt.wantBallX)
			}
			if absFloat(game.BallY-tt.wantBallY) > epsilon {
				t.Errorf("calcMoveBall() BallY = %v, want %v", game.BallY, tt.wantBallY)
			}
		})
	}
}
