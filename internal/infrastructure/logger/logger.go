package logger

import (
	"log/slog"
	"os"

	"github.com/valenoirs/flashcard-api/internal/config"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if cfg.App.Env == "development" {
		opts.Level = slog.LevelDebug
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
