package config

import (
	"errors"
	"fmt"
	"log/slog"
)

type getEnv func(string) string

// Load reads and validates configuration using a two-phase process:
// 1. Loads all environment variables and aggregates any parsing errors.
// 2. Validates the configuration against business logic rules.
func Load(getenv getEnv) (*Config, error) {
	r := newReader(getenv)

	cfg := &Config{
		App:    loadApp(r),
		Http:   loadHTTP(r),
		Logger: loadLogger(r),
	}

	// Return any errors encountered while reading environment variables.
	if err := r.err(); err != nil {
		return nil, fmt.Errorf("load config:\n%w", err)
	}

	var errs []error

	// Run business logic validation on the loaded configuration.
	if err := cfg.validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}

// loadApp loads application-level configuration
func loadApp(r *reader) *APP {
	return &APP{
		Envrionment: required(
			r,
			envAppEnrionment,
			oneOf("production", "staging", "development"),
		),
	}
}

// loadHTTP loads HTTP server configuration
func loadHTTP(r *reader) *HTTP {
	return &HTTP{
		Addr: optional(
			r,
			envHTTPAddr,
			":8080",
			asString,
		),
	}
}

// loadLogger loads Logger configuration
func loadLogger(r *reader) *Logger {
	return &Logger{
		LOG_LEVEL: optional(
			r,
			envLogLevel,
			slog.LevelDebug,
			asLevel,
		),
		LOG_FORMAT: optional(
			r,
			envLogFormat,
			FormatText,
			oneOf(FormatJSON, FormatText),
		),
		AddSource: optional(
			r,
			envAddSource,
			false,
			asBool,
		),
	}
}
