// Package middleware contains the HTTP middleware chain.
//
// Composition order in router.go (outer to inner):
//  1. RequestID    - generate X-Request-Id if absent, store in context
//  2. Logging      - structured access log with duration and request ID
//  3. Recovery     - catch panics, log them, return 500
//  4. CORS         - permissive for the TUI client (no browser caller for now)
//  5. Auth         - validate JWT and attach user to context (mounted on
//     protected routes only)
// TODO(phase-5): Auth (JWT validation, user context injection)

package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-Id"
	// maxRequestIDLen bounds an inbound ID so an attacker-controlled header
	// can't bloat every log line for the request.
	maxRequestIDLen = 128
)

// RequestID is middleware that generates X-Request-Id if absent, stores into
// request context. Outermost middleware to ensure all requests have an ID so
// they can be tracked in the code.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if !validRequestID(id) { // honour a sane inbound ID; else mint one
			id = newRequestID()
		}
		w.Header().Set(requestIDHeader, id)   // echo to the client
		ctx := WithRequestID(r.Context(), id) // AND stash for our own code
		next.ServeHTTP(w, r.WithContext(ctx)) // pass the NEW request down
	})
}

// validRequestID reports whether an inbound ID is safe to echo and log: a
// non-empty, bounded string of printable ASCII. Anything else is replaced
// with a freshly generated ID.
func validRequestID(s string) bool {
	// valid length
	if s == "" || len(s) > maxRequestIDLen {
		return false
	}
	// valid characters
	for _, c := range s {
		if c < 0x20 || c > 0x7e {
			return false
		}
	}
	return true
}

// newRequestID returns a random UUID string, assumed to be unique
func newRequestID() string {
	return uuid.NewString()
}

// Logging middleware constructs a wrapper to record the status of actioned
// requests, such that they can be logged. Next inner from RequestID so that
// all requests are logged.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r) // pass the status recorder, not w

		// log request
		slog.Info("http request",
			"request_id", RequestIDFromContext(r.Context()), // reads what RequestID stashed
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"bytes", rec.bytes,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote", r.RemoteAddr,
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter // inherit Header(), Write(), WriteHeader()
	status              int
	bytes               int
	wrote               bool
}

// Overwrite WriteHeader() so status is recorded
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.wrote = true
	rec.ResponseWriter.WriteHeader(code)
}

// Overwrite Write() so write state recorded
func (rec *statusRecorder) Write(b []byte) (int, error) {
	rec.wrote = true
	n, err := rec.ResponseWriter.Write(b)
	rec.bytes += n
	return n, err
}

// Recovery middleware handles panics, logging and returning error info.
// If no panic, continue. If panic:
// - Check if handler intentionally aborted, re-panic
// - Log error & write problem
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return // normal path: nothing panicked
			}
			if rec == http.ErrAbortHandler {
				panic(rec) // intentional abort: re-panic, don't swallow
			}
			slog.Error("panic recovered",
				"request_id", RequestIDFromContext(r.Context()),
				"panic", rec,
				"path", r.URL.Path,
			)
			if !responseStarted(w) {
				WriteProblem(w, r, http.StatusInternalServerError, "")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// responseStarted checks if ResponseWrite has wrote yet
func responseStarted(w http.ResponseWriter) bool {
	rec, ok := w.(*statusRecorder)
	return ok && rec.wrote
}
