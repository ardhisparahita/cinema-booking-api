package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return slog.New(handler)
}
