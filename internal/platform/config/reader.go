package config

import (
	"errors"
	"fmt"
	"strings"
)

// reader reads environment variables and collects any parsing errors,
// allowing startup to report all bad settings at once instead of failing one-by-one.
type reader struct {
	getenv getEnv
	errs   []error
}

func newReader(getenv getEnv) *reader {
	return &reader{getenv: getenv}
}

// err combines all collected reading errors into a single error.
func (r *reader) err() error {
	return errors.Join(r.errs...)
}

// fail records a parsing or configuration failure for a specific key.
func (r *reader) fail(key string, cause error) {
	r.errs = append(r.errs, fmt.Errorf("%s: %w", key, cause))
}

// optional reads an environment variable,
// returning the parsed value or a default if empty/invalid.
func optional[T any](r *reader, key string, def T, parse func(string) (T, error)) T {
	raw := strings.TrimSpace(r.getenv(key))

	if raw == "" {
		return def
	}

	v, err := parse(raw)

	if err != nil {
		r.fail(key, err)
		return def
	}

	return v
}

// required reads and parses a mandatory environment variable,
// recording an error if missing or invalid.
func required[T any](r *reader, key string, parse func(string) (T, error)) T {
	var zero T

	raw := strings.TrimSpace(r.getenv(key))

	if raw == "" {
		r.fail(key, errors.New("required but not set"))
		return zero
	}

	v, err := parse(raw)
	if err != nil {
		r.fail(key, err)
		return zero
	}

	return v
}
