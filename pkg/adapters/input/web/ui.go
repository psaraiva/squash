package web

import (
	"fmt"

	"github.com/psaraiva/squash/internal/app"
	"github.com/psaraiva/squash/internal/ports"
)

func PaintGame(r ports.Renderer, p *app.Squash) {
	r.Clear()

	drawTextScore(r, getTextScore(p))
	drawTextLives(r, p, getTextLives(p))

	switch p.State {
	case app.StateMenu:
		drawTextCenter(r, p, getTextStateMenu())

	case app.StatePaused:
		drawTextCenter(r, p, getTextStatePaused())

	case app.StatePlaying:
		drawGameElements(r, p)

	case app.StateGameOver:
		drawTextCenter(r, p, getTextStateGameOver(p))
	}

	if p.DebugMode {
		drawDebugInfo(r, getDebugInfo(p))
	}
}

func getTextScore(p *app.Squash) string {
	return fmt.Sprintf("Score: %d", p.Score)
}

func drawTextScore(r ports.Renderer, text string) {
	r.DrawText(text, 30, 30)
}

func getTextLives(p *app.Squash) string {
	return fmt.Sprintf("Lives: %d", p.Lives)
}

func drawTextLives(r ports.Renderer, p *app.Squash, text string) {
	r.DrawText(text, p.Width-120, 30)
}

func drawGameElements(r ports.Renderer, p *app.Squash) {
	r.DrawBall(p.BallX, p.BallY, p.BallSize)
	r.DrawPaddle(p.PaddleX, p.PaddleY, p.PaddleW, p.PaddleH)
}

func getDebugInfo(p *app.Squash) []string {
	return []string{
		fmt.Sprint("Game:....."),
		fmt.Sprintf("FPS:      %d", p.Fps),
		fmt.Sprintf("Level:    %d", p.LastLevel),
		fmt.Sprint("Ball:....."),
		fmt.Sprintf("Size:     %.1f", p.BallSize),
		fmt.Sprintf("Spawn:    [%.1f, %.1f]", p.BallSpawnX, p.BallSpawnY),
		fmt.Sprintf("Position: [%.1f, %.1f]", p.BallX, p.BallY),
		fmt.Sprintf("Velocity: [%.2f, %.2f]", p.BallDX, p.BallDY),
	}
}

func drawDebugInfo(r ports.Renderer, info []string) {
	for i, text := range info {
		r.DrawDebugText(text, 10, 80+float64(i*15))
	}
}

func getTextStateMenu() []string {
	return []string{
		"SQUASH - LEFT CLICK TO START",
		"(RIGHT CLICK TO PAUSE)",
	}
}

func getTextStatePaused() []string {
	return []string{"PAUSED - RIGHT CLICK TO RESUME"}
}

func getTextStateGameOver(p *app.Squash) []string {
	return []string{
		fmt.Sprintf("GAME OVER - SCORE: %d", p.Score),
		"(LEFT CLICK TO RESTART)",
	}
}

func drawTextCenter(r ports.Renderer, p *app.Squash, text []string) {
	for i, line := range text {
		r.DrawText(line, (p.Width-r.MeasureText(line))/2, p.Height/2+float64(i*20))
	}
}
