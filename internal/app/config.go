package app

type Config struct {
	Debug          bool
	InitialLives   int
	InitialLevel   int
	SpeedIncrement float64
	BallScale      float64
	Fps            int
	DeltaTime      float64
}

func NewDefaultConfig() Config {
	return Config{
		Debug:          false,
		InitialLives:   3,
		InitialLevel:   0,
		SpeedIncrement: 0.25,
		BallScale:      0.0,
		Fps:            30,
	}
}
