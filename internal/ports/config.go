package ports

import game "github.com/psaraiva/squash/internal/app"

type ConfigProvider interface {
	Load() game.Config
}
