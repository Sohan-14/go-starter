package middleware

import (
	"log/slog"
	"net/http"
	"slices"
)

// Middleware is a function that wraps an http.Handler to add behaviour.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to a handler in order.
// The first middleware in the list is the outermost (executes first on request,
// last on response).
func chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		h = middleware(h)
	}
	return h
}

// Standard returns the default middleware chain for all HTTP requests.
// This is the single source of truth for middleware order and composition.
// Add new global middlewares here, not in main.go.
func Standard(log *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return chain(h,
			panicRecover(log),
			requestID(),
			accessLog(log),
			// drain(d),
		)

	}
}
