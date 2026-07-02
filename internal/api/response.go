package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Error responses follow RFC 7807 (application/problem+json). A single
// envelope keeps error shapes consistent across every handler, and the
// mapping from domain errors to HTTP status codes lives in WriteError alone,
// so handlers stay thin and never hand-roll an error body.
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
