package middleware

import (
	"log/slog"
	"net/http"

	"github.com/valenoirs/flashcard-api/internal/config"
)

type Manager struct {
	CORS func(http.Handler) http.Handler
}

func NewManager(cfg *config.Config, logger *slog.Logger) *Manager {
	return &Manager{
		CORS: CORS(cfg, logger),
	}
}
