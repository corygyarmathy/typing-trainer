// Package api wires the HTTP router and middleware chain.
//
// Each bounded context (auth, progress, session) exposes a RegisterRoutes
// function that takes a Router. This package composes them and applies
// shared middleware (request ID, logging, recovery, auth).
package api

import (
	"context"
	"net/http"
	"time"
)

// chain avoids having to nest each middleware handler, increases readability.
//
// Reverses order so that the first middleware is the outermost (runs first)
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- { // apply middlerware in reverse
		h = mws[i](h)
	}
	return h
}

// Router constructs the application HTTP handler.
//
// Dependencies (services) are passed in here, not pulled from globals.
func Router(ready func(ctx context.Context) error) http.Handler {
	// Initiate multiplexer / router
	mux := http.NewServeMux()
	// Register handler functions (patterns that the API can handle)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /readyz", readyz(ready))
	// Construct and return Handler, composing middleware handlers
	return chain(mux, RequestID, Logging, Recovery) // Order matters: outer-> inner
}

// healthz returns status ok; basic check that API JSON response writer is functioning
func healthz(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// readyz confirms database is reachable; returning status ok if db ping
// succeeds within 2 seconds
func readyz(ready func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// attempt to ping db; return status service unavailable if timed out
		if err := ready(ctx); err != nil {
			WriteProblem(w, r, http.StatusServiceUnavailable, "database is not reachable")
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	}
}
