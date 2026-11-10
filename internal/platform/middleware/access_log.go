package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// accessLog logs every HTTP request with method, path, status, duration,
// and bytes written. It uses slog.InfoContext so the request_id from
// the context is automatically included.
func accessLog(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				responseWriter := newResponseWriter(w)

				next.ServeHTTP(responseWriter, r)

				log.InfoContext(r.Context(), "http request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", responseWriter.status),
					slog.Duration("duration", time.Since(start)),
					slog.Int64("bytes", responseWriter.bytes),
					slog.String("remote_addr", r.RemoteAddr),
				)
			},
		)
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
// and the number of bytes written.
type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	if code >= 100 && code < 200 && code != http.StatusSwitchingProtocols {
		rw.ResponseWriter.WriteHeader(code)
		return
	}

	rw.wroteHeader = true
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.wroteHeader = true
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += int64(n)
	return n, err
}

func (rw *responseWriter) Written() bool {
	return rw.wroteHeader
}

func (rw *responseWriter) Flush() {
	rw.wroteHeader = true
	_ = http.NewResponseController(rw.ResponseWriter).Flush()
}
