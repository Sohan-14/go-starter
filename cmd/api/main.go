package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sohan-14/go-starter/internal/platform/config"
	"github.com/sohan-14/go-starter/internal/platform/database"
	"github.com/sohan-14/go-starter/internal/platform/logger"
	"github.com/sohan-14/go-starter/internal/platform/middleware"
	"github.com/sohan-14/go-starter/internal/platform/server"
)

func main() {
	cfg, err := config.Load(os.Getenv)

	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	log := logger.Setup(cfg.Logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, log); err != nil {
		log.Error("application stopped with error", "error", err)
		os.Exit(1)
	}

	log.Info("application stopped gracefully")
}

func run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	db, err := database.New(ctx, cfg.Database, log, "go-starter")
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer db.Close()

	srv := server.New(cfg.HTTP, log)
	registerRoutes(srv)

	srv.ApplyMiddleware(middleware.Standard(log))

	log.InfoContext(ctx, "application starting", "environment", cfg.App.Envrionment)

	if err := srv.Run(ctx); err != nil {
		log.Error("application stopped with error", "error", err)
		os.Exit(1)
	}

	return nil
}

// registerRoutes attaches infrastructure and health check endpoints to the server.
// This keeps transport concerns (server) separate from routing concerns (main).
func registerRoutes(srv *server.Server) {
	srv.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if srv.IsDraining() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if srv.IsDraining() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// TODO: Add database, cache, and downstream dependency checks here later.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})
}
