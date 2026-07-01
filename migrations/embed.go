// Package migrations embeds the goose SQL migration files so they travel with
// the compiled binary and can be applied at startup, with no migrations
// directory needed on the deployment host.
//
// A go;embed directive cannot reference paths outside its own package
// directory, so the embed declaration lives here, in the migrations directory
// itself, rather than in internal/platform/database.
package migrations

import "embed"

// SQLMigrationsFS holds every goose SQL migration in this directory.
//
//go:embed *.sql
var SQLMigrationsFS embed.FS
