package ports

type Renderer interface {
	Clear()
	DrawBall(x, y, radius float64)
	DrawDebugText(text string, x, y float64)
	DrawPaddle(x, y, w, h float64)
	DrawText(text string, x, y float64)
	MeasureText(text string) float64
}
