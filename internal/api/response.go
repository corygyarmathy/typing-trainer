// Package api: shared JSON response and error helpers.
package api

// TODO: Decide the API error format and write it down — a consistent error envelope
// or RFC 7807 problem+json. There's `response.go` for shared helpers, which is the
// right instinct, but no documented decision, and inconsistent error shapes are
// something easily noticed in thirty seconds of curling.

// TODO(phase-4): implement
//   - WriteJSON(w http.ResponseWriter, status int, v any)
//   - WriteError(w http.ResponseWriter, r *http.Request, err error)
//
// The error helper should map domain errors (e.g. errors.Is(err, ErrNotFound))
// to HTTP status codes in one place, so handlers stay thin.
