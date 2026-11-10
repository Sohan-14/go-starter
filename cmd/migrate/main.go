package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sohan-14/go-starter/internal/platform/config"
	"github.com/sohan-14/go-starter/internal/platform/logger"
	"github.com/sohan-14/go-starter/internal/platform/migration"
	"github.com/sohan-14/go-starter/migrations"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	log := logger.Setup(cfg.Logger)

	runner, err := migration.New(cfg.Database.URL, migrations.FS, log)
	if err != nil {
		log.Error("failed to create migration runner", "error", err)
		os.Exit(1)
	}
	defer runner.Close()

	command := os.Args[1]

	switch command {
	case "up":
		if err := runner.Up(); err != nil {
			log.Error("migration up failed", "error", err)
			os.Exit(1)
		}
	case "down":
		if err := runner.Down(); err != nil {
			log.Error("migration down failed", "error", err)
			os.Exit(1)
		}
	case "status":
		if err := runner.Status(); err != nil {
			log.Error("migration status failed", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}
}
