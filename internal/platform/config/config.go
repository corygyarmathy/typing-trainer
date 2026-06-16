// Package config loads runtime configuration from environment variables.
//
// Follows 12-factor principles: configuration is environment, secrets never
// live in source, and the zero-config defaults assume local development.
package config

// Config holds all runtime configuration for the server.
//
// TODO(phase-1): populate fields as they become needed:
//   - HTTPAddr      string
//   - DatabaseURL   string
//   - JWTSecret     string
//   - LogLevel      string
//   - MetricsAddr   string
type Config struct{}

// Load reads configuration from the process environment.
// Returns an error if any required value is missing or invalid.
func Load() (Config, error) {
	return Config{}, nil
}
