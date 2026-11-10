package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/sohan-14/go-starter/internal/platform/logger"
)

const headerRequestID = "X-Request-ID"

// requestId generates a unique request ID or extracts one from the incoming
// X-Request-ID header. It injects the ID into the request context via the
// logger's ContextWithAttrs and sets it on the response header.
func requestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				id := r.Header.Get(headerRequestID)

				if id == "" {
					id = generateRequestID()
				}

				w.Header().Set(headerRequestID, id)

				ctx := logger.ContextWithAttrs(r.Context(),
					slog.String("request_id", id),
				)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
