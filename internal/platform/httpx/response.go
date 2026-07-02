// Package httpx holds the shared HTTP kit used across every bounded context:
// JSON and RFC 7807 problem+json response writers, request-scoped context
// helpers, and the cross-cutting middleware. It imports no domain package, so
// domains can depend on it without forming an import cycle; route composition
// lives above it in cmd/server.
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Error responses follow RFC 7807 (application/problem+json). Routing every
// handler error through WriteProblem keeps error bodies consistent and stops
// handlers hand-rolling their own JSON. The caller passes the HTTP status and
// a safe, human-readable detail explicitly; there is no domain-error-to-status
// mapping yet (that arrives with the service layer in Phase 4).
//
// Note: 404 and 405 are emitted by http.ServeMux before any handler runs, so
// they bypass this writer and carry the stdlib's text/plain body. See ADR 0019.
//
// Refer: https://datatracker.ietf.org/doc/html/rfc7807

// Problem is an RFC 7807 problem detail. Only the fields v1 uses are modelled;
// the spec allows extension members to be added later without breaking clients.
type Problem struct {
	// Type is a URI reference identifying the problem type. "about:blank" is
	// the RFC default and means "see the status code".
	Type string `json:"type"`
	// Title is a short, human-readable summary, stable for a given Type.
	Title string `json:"title"`
	// Status is the HTTP status code, duplicated in the body per the RFC.
	Status int `json:"status"`
	// Detail is a human-readable explanation specific to this occurrence.
	// Don't pass err. Prevents internal server logic from being leaked extenally.
	Detail string `json:"detail,omitempty"`
	// Instance identifies the specific occurrence using the request ID.
	Instance string `json:"instance,omitempty"`
}

// WriteJSON encodes v as JSON with the given status code. It is the single
// place response headers and encoding are handled for success responses.
func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Headers are already flushed; the most we can do is record it.
		slog.Error("encoding JSON response", "err", err)
	}
}

// WriteProblem writes an RFC 7807 problem+json response. detail may be empty,
// in which case clients fall back to the status code's standard meaning.
func WriteProblem(w http.ResponseWriter, r *http.Request, statusCode int, detail string) {
	p := Problem{
		Type:     "about:blank",
		Title:    http.StatusText(statusCode),
		Status:   statusCode,
		Detail:   detail,
		Instance: RequestIDFromContext(r.Context()),
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("encoding problem response", "err", err)
	}
}
