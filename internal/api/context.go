package api

import "context"

// ctxKey is unexported so no other package can collide with our context keys.
type ctxKey int

const (
	requestIDKey ctxKey = iota
	userIDKey
)

// WithRequestID returns a copy of ctx carrying the given request ID. The
// middleware sets it; handlers and the problem-detail encoder read it back.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFromContext returns the request ID stored by the middleware, or
// the empty string if none is present.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
