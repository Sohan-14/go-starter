package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sohan-14/go-starter/internal/platform/config"
)

// Server wraps http.Server with graceful shutdown support.
// It owns the lifecycle: creation, startup, drain, and stop.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	httpConfig *config.HTTP
	log        *slog.Logger
	draining   atomic.Bool
}

// New creates a configured Server.
func New(httpConfig *config.HTTP, log *slog.Logger) *Server {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:              httpConfig.Addr,
		Handler:           mux,
		ReadHeaderTimeout: httpConfig.ReadHeaderTimeout,
		ReadTimeout:       httpConfig.ReadTimeout,
		WriteTimeout:      httpConfig.WriteTimeout,
		IdleTimeout:       httpConfig.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		httpConfig: httpConfig,
		log:        log,
	}
}

// Handle registers an http.Handler for the given pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers an http.HandlerFunc for the given pattern.
func (s *Server) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handlerFunc)
}

// IsDraining returns true if the server is in the process of shutting down.
// Middleware can use this to return 503 Service Unavailable for new requests.
func (s *Server) IsDraining() bool {
	return s.draining.Load()
}

// Run binds the TCP port and starts the server.
// Binding synchronously ensures we fail fast if the port is already in use.
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.httpServer.Addr, err)
	}
	return s.Serve(ctx, ln)
}

// Serve starts the HTTP server on the provided listener and blocks until
// the context is canceled or the server crashes.
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- s.httpServer.Serve(ln)
	}()

	s.log.InfoContext(ctx, "http server listening", "addr", ln.Addr().String())

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	return s.shutdown(ctx)
}

// shutdown gracefully drains in-flight requests.
func (s *Server) shutdown(ctx context.Context) error {
	s.draining.Store(true)

	// Create a detached context for shutdown operations so we don't log
	// with a canceled context, which can confuse downstream middleware.
	detachedCtx := context.WithoutCancel(ctx)

	s.log.InfoContext(detachedCtx, "shutting down",
		"drain_delay", s.httpConfig.DrainDelay.String(),
		"timeout", s.httpConfig.ShutdownTimeout.String(),
	)

	// Drain delay allows load balancers to detect the pod is terminating
	// and stop sending new traffic before we actually close connections.
	if s.httpConfig.DrainDelay > 0 {
		time.Sleep(s.httpConfig.DrainDelay)
	}

	shutdownCtx, cancel := context.WithTimeout(detachedCtx, s.httpConfig.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		// If graceful shutdown fails or times out, force close immediately.
		_ = s.httpServer.Close()
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	s.log.InfoContext(detachedCtx, "http server stopped")
	return nil
}
