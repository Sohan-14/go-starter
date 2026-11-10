package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// panicRecover catches panics from any downstream handler or middleware,
// logs the full stack trace, and returns a 500 response.
// This prevents a single bad request from crashing the entire server.
func panicRecover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						if err == http.ErrAbortHandler {
							panic(err)
						}

						log.ErrorContext(r.Context(), "panic recovered",
							slog.Any("error", err),
							slog.String("stack", string(debug.Stack())),
						)

						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}()

				next.ServeHTTP(w, r)
			},
		)
	}
}
