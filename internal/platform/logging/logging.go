// Package logging configures the process-wide slog logger.
//
// Output is JSON in production, text in development, selected by the
// APP_ENV environment variable. The configured logger is installed as
// slog.Default so packages that don't take a logger dependency still
// produce structured output.
package logging

import "log/slog"

// Setup installs a slog logger as the default based on the level string.
// Valid levels: "debug", "info", "warn", "error".
//
// TODO(phase-1):
//   - Use slog.NewJSONHandler for production, slog.NewTextHandler for dev
//   - Add a "service" attribute to every log line
//   - Return the logger so callers that prefer DI can use it
func Setup(level string) *slog.Logger {
	return slog.Default()
}
