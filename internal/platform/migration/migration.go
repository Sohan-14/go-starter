package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

type Runner struct {
	db  *sql.DB
	log *slog.Logger
}

func New(dbURL string, migrationsFS fs.FS, log *slog.Logger) (*Runner, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set dialect: %w", err)
	}

	goose.SetLogger(goose.NopLogger())

	return &Runner{
		db:  db,
		log: log,
	}, nil
}

func (r *Runner) Up() error {
	if err := goose.Up(r.db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	r.log.Info("migrations applied successfully")
	return nil
}

func (r *Runner) Down() error {
	if err := goose.Down(r.db, "."); err != nil {
		return fmt.Errorf("rollback migration: %w", err)
	}
	r.log.Info("migration rolled back successfully")
	return nil
}

func (r *Runner) Status() error {
	if err := goose.Status(r.db, "."); err != nil {
		return fmt.Errorf("migration status: %w", err)
	}
	return nil
}

func (r *Runner) Close() error {
	return r.db.Close()
}
