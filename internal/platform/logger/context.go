package logger

import (
	"context"
	"log/slog"
)

// ctxKey is a private type used to uniquely
// store and retrieve log attributes from a context.
type ctxKey struct{}

// ContextWithAttrs injects additional slog attributes into the provided context,
// merging them with any existing attributes already stored in the context.
func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing, _ := ctx.Value(ctxKey{}).([]slog.Attr)

	merged := make([]slog.Attr, 0, len(existing)+len(attrs))
	merged = append(merged, existing...)
	merged = append(merged, attrs...)

	return context.WithValue(ctx, ctxKey{}, merged)
}

// contextHandler is a custom slog.Handler wrapper that automatically extracts
// context-bound attributes and attaches them to log records during handling.
type contextHandler struct {
	slog.Handler
}

// Handle intercepts log records, appends any attributes found in the context,
// and delegates handling to the underlying slog.Handler.
func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(ctxKey{}).([]slog.Attr); ok && len(attrs) > 0 {
		r = r.Clone()
		r.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new contextHandler with the given
// attributes added to the handler's underlying state.
func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{h.Handler.WithAttrs(attrs)}
}

// WithGroup returns a new contextHandler with the given
// group name added to the handler's underlying state.
func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{h.Handler.WithGroup(name)}
}
