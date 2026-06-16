// Package database provides the pgx connection pool and migration runner.
//
// The package exposes a minimal surface: Open returns a *pgxpool.Pool wired
// to the configured database, and Migrate runs all pending goose migrations
// from the embedded migrations directory.
//
// Repositories elsewhere in the codebase accept *pgxpool.Pool (or the
// narrower pgx.Tx for transactional work) rather than this package directly,
// so they remain testable with testcontainers.
package database

import "context"

// Open constructs a database pool from the given DSN.
//
// TODO(phase-2):
//   - Return *pgxpool.Pool
//   - Set sane defaults: MaxConns, MinConns, HealthCheckPeriod
//   - Verify connectivity via pool.Ping before returning
func Open(ctx context.Context, dsn string) (any, error) {
	return nil, nil
}

// Migrate applies all pending migrations from the embedded migrations dir.
//
// TODO(phase-2):
//   - Use //go:embed to bundle ../../../migrations/*.sql
//   - Acquire a connection from the pool and pass to goose.Up
//   - Take an advisory lock so concurrent boots don't race
func Migrate(ctx context.Context, pool any) error {
	return nil
}
