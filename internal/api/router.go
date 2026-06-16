// Package api wires the HTTP router and middleware chain.
//
// Each bounded context (auth, progress, session) exposes a RegisterRoutes
// function that takes a chi.Router. This package composes them and applies
// shared middleware (request ID, logging, recovery, auth).
package api

import "net/http"

// Router constructs the application HTTP handler.
//
// TODO(phase-4):
//   - chi.NewRouter()
//   - apply middleware chain (recovery, requestid, logging, cors)
//   - mount /api/v1 subrouter and call each context's RegisterRoutes
//   - mount /healthz, /readyz, /metrics
//
// Dependencies (services) are passed in here, not pulled from globals.
func Router( /* services... */ ) http.Handler {
	return http.NotFoundHandler()
}
