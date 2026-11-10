package middleware

import "net/http"

// drainer is an interface for checking if the server is draining.
// This avoids a circular dependency between the middleware and server packages.
type drainer interface {
	isDraining() bool
}

// draib rejects new requests with 503 Service Unavailable when the server
// is in the process of shutting down. This allows load balancers to stop
// routing traffic while in-flight requests complete.
func drain(d drainer) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if d.isDraining() {
					http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
