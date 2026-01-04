package app

import (
	"math/rand"
)

const (
	BaseSpeedBall      float64 = 200.0
	PointsPerLevel     int     = 100
	PointsPerCollision int     = 10
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

type Squash struct {
	// General
	DeltaTime      float64
	Fps            int
	Height         float64
	LastLevel      int
	Lives          int
	Score          int
	State          GameState
	SpeedIncrement float64
	Width          float64

	// elements
	BallSize               float64
	BallX, BallY           float64
	BallDX, BallDY         float64 // direction
	BallSpawnX, BallSpawnY float64
	PaddleX, PaddleY       float64
	PaddleW, PaddleH       float64

	DebugMode bool
}

func NewSquash(w, h float64, cfg Config) *Squash {
	p := &Squash{
		Width:   w,
		Height:  h,
		PaddleH: 60,
		PaddleW: 10,
	}

	p.PaddleX = 10

	p.Reset(cfg)
	p.State = StateMenu
	return p
}

func (p *Squash) Reset(cfg Config) {
	p.loadConfigDefaul(cfg)
	p.BallSize = calcBallSize(cfg.BallScale)
	p.PaddleY = calcRespawPaddleY(p.Height, p.PaddleH)
	p.respawnBall()
}

func (p *Squash) loadConfigDefaul(cfg Config) {
	p.DebugMode = cfg.Debug
	p.Fps = cfg.Fps
	p.DeltaTime = cfg.DeltaTime
	p.Lives = cfg.InitialLives
	p.LastLevel = cfg.InitialLevel
	p.SpeedIncrement = cfg.SpeedIncrement
	p.Score = cfg.InitialLevel * 100
}

func (p *Squash) respawnBall() {
	p.BallX = p.Width / 2
	p.BallY = calcBallStartY(p.Height, rand.Float64())

	p.BallSpawnX = p.BallX
	p.BallSpawnY = p.BallY

	factor := calcSpeedFactor(p.LastLevel, p.SpeedIncrement)
	p.BallDX = BaseSpeedBall * factor
	p.BallDY = BaseSpeedBall * factor

	p.BallDX, p.BallDY = calcRandomDirectionStartBall(p.BallDX, p.BallDY)
}

func calcSpeedFactor(level int, increment float64) float64 {
	if level < 0 || level > 99 {
		level = 0
	}

	if increment < 0 || increment > 1.0 {
		return 1.0
	}

	return 1.0 + (float64(level) * increment)
}

func calcBallSize(scale float64) float64 {
	baseSize := 10.0
	if scale < 0 || scale > 1.0 {
		return baseSize
	}

	return baseSize * (1.0 + scale)
}

func calcBallStartY(height float64, randomSource float64) float64 {
	margin := height * 0.15
	playable := height * 0.70
	return margin + (randomSource * playable)
}

func calcRespawPaddleY(height, paddleH float64) float64 {
	return (height / 2) - (paddleH / 2)
}

func calcRandomDirectionStartBall(ballDX, ballDY float64) (float64, float64) {
	if rand.Float64() > 0.5 {
		ballDX *= -1
	}

	if rand.Float64() > 0.5 {
		ballDY *= -1
	}

	return ballDX, ballDY
}
