package database

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sohan-14/go-starter/internal/platform/config"
)

type DB struct {
	pool             *pgxpool.Pool
	log              *slog.Logger
	readinessTimeout time.Duration
}

func New(ctx context.Context, dbConfig *config.Database, log *slog.Logger, appName string) (*DB, error) {
	poolCfg, err := poolConfig(dbConfig, appName)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingCtx, cencel := context.WithTimeout(ctx, dbConfig.StartupTimeout)
	defer cencel()

	if err := waitReady(pingCtx, pool, log); err != nil {
		pool.Close()
		return nil, err
	}

	log.InfoContext(ctx, "database connected",
		"host", poolCfg.ConnConfig.Host,
		"port", poolCfg.ConnConfig.Port,
		"database", poolCfg.ConnConfig.Database,
		"max_conns", poolCfg.MaxConns,
		"min_conns", poolCfg.MinConns,
	)

	return &DB{
		pool:             pool,
		log:              log,
		readinessTimeout: dbConfig.ReadinessTimeout,
	}, nil
}

// Pool exposes the pgx pool.
// sqlc-generated repositories will use this later.
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// Health checks database readiness with a bounded timeout.
// Use this from /readyz, not from /healthz.
func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, db.readinessTimeout)
	defer cancel()

	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check: %w", err)
	}

	return nil
}

// Close closes all database connections.
// Call this after the HTTP server has stopped.
func (db *DB) Close() {
	db.log.Info("closing database connection pool")
	db.pool.Close()
}

// poolConfig build the *pgxpool.Config using environment DB config
func poolConfig(dbConfig *config.Database, appName string) (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(dbConfig.URL)

	if err != nil {
		return nil, fmt.Errorf("database parse url: %w", err)
	}

	poolCfg.MaxConns = int32(dbConfig.MaxConns)
	poolCfg.MinConns = int32(dbConfig.MinConns)
	poolCfg.MaxConnLifetime = dbConfig.MaxConnLifetime
	poolCfg.MaxConnIdleTime = dbConfig.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = dbConfig.HealthCheckPeriod
	poolCfg.ConnConfig.ConnectTimeout = dbConfig.ConnectTimeout

	poolCfg.MaxConnLifetimeJitter = dbConfig.MaxConnLifetime / 10

	poolCfg.ConnConfig.RuntimeParams["application_name"] = appName
	poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = strconv.FormatInt(dbConfig.StatementTimeout.Milliseconds(), 10)

	return poolCfg, nil
}

// waitReady mulitple attemps for ping to the database
// until the startup timeout set in the environment db config
func waitReady(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger) error {
	delay := 250 * time.Millisecond

	for attempt := 1; ; attempt++ {
		err := pool.Ping(ctx)
		if err == nil {
			return nil
		}

		log.WarnContext(ctx, "database not ready", "attempt", attempt)

		select {
		case <-ctx.Done():
			return fmt.Errorf("database not reachable after %d attempts", attempt)
		case <-time.After(delay):
		}

		delay = min(delay*2, 2*time.Second)
	}
}
