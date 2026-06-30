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

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Open constructs a database pool from the given connection string.
//
// TODO(phase-2):
//   - Return *pgxpool.Pool
//   - Set sane defaults: MaxConns, MinConns, HealthCheckPeriod
//   - Verify connectivity via pool.Ping before returning
func Open(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return &pgxpool.Pool{}, fmt.Errorf("failed to parse db connString config: %w", err)
	}

	cfg.MaxConns = 4
	cfg.MinConns = 1
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return &pgxpool.Pool{}, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return &pgxpool.Pool{}, fmt.Errorf("pinging database: %w", err)
	}

	return pgxpool.New(ctx, connString)
}

// Migrate applies all pending migrations from the embedded migrations dir.
//
// TODO(phase-2):
//   - Use //go:embed to bundle ../../../migrations/*.sql
//   - Acquire a connection from the pool and pass to goose.Up
//   - Take an advisory lock so concurrent boots don't race
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	return nil
}
