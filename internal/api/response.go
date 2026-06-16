// Package api: shared JSON response and error helpers.
package api

// TODO(phase-4): implement
//   - WriteJSON(w http.ResponseWriter, status int, v any)
//   - WriteError(w http.ResponseWriter, r *http.Request, err error)
//
// The error helper should map domain errors (e.g. errors.Is(err, ErrNotFound))
// to HTTP status codes in one place, so handlers stay thin.
