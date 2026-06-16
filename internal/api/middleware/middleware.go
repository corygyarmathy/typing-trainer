// Package middleware contains the HTTP middleware chain.
//
// Composition order in router.go (outer to inner):
//  1. RequestID    - generate X-Request-Id if absent, store in context
//  2. Logging      - structured access log with duration and request ID
//  3. Recovery     - catch panics, log them, return 500
//  4. CORS         - permissive for the TUI client (no browser caller for now)
//  5. Auth         - validate JWT and attach user to context (mounted on
//     protected routes only)
package middleware

// TODO(phase-4): RequestID, Logging, Recovery
// TODO(phase-5): Auth (JWT validation, user context injection)
