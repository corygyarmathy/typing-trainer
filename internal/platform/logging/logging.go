// Package logging configures the process-wide slog logger.
//
// Output is JSON in production, text in development, selected by the
// APP_ENV environment variable. The configured logger is installed as
// slog.Default so packages that don't take a logger dependency still
// produce structured output.
package logging

import (
	"log/slog"
	"os"
)

// Setup installs the process-wide default logger and returns it for DI callers.
// Valid levels: "debug", "info", "warn", "error".
func Setup(level, env string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo // config already validated; safety net
	}

	opts := &slog.HandlerOptions{Level: lvl}
	var h slog.Handler
	if env == "production" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(h).With("service", "typist")
	slog.SetDefault(logger)
	return logger
}
