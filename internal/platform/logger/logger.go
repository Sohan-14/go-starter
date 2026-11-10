package logger

import (
	"log/slog"
	"os"

	"github.com/sohan-14/go-starter/internal/platform/config"
)

// new creates a new slog.Logger based on the provided config.
func new(logger *config.Logger) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: logger.AddSource,
		Level:     logger.Level,
	}

	var handler slog.Handler

	switch logger.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(contextHandler{handler})
}

// Setup creates a new logger and sets it as the default slog logger.
func Setup(logger *config.Logger) *slog.Logger {
	l := new(logger)
	slog.SetDefault(l)
	return l
}
