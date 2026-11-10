package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const maxPoolConns = 1000

// validate runs business logic validation across all configuration sections.
func (c *Config) validate() error {
	var errs []error

	if err := c.App.validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HTTP.validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Logger.validate(c.App.Envrionment); err != nil {
		errs = append(errs, err)
	}

	if err := c.Database.validate(c.App.Envrionment); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// validate ensures the application level validation and business logic
func (a *APP) validate() error {
	if a.Envrionment != "development" && a.Envrionment != "production" && a.Envrionment != "stagging" {
		return fmt.Errorf("%q must be one of these development, production, stagging", envAppEnrionment)
	}
	return nil
}

// validate ensures the http server level validation and business logic
func (h *HTTP) validate() error {
	_, port, err := net.SplitHostPort(h.Addr)

	if err != nil {
		return fmt.Errorf("%s: want host:port such as :8080: %w", envHTTPAddr, err)
	}

	n, err := strconv.Atoi(port)
	if err != nil || n < 0 || n > 65535 {
		return fmt.Errorf("%s: port %q must be a number from 0 to 65535", envHTTPAddr, port)
	}

	return nil
}

// validate ensures logger configuration complies with environment-specific business rules.
func (l *Logger) validate(environment string) error {
	// In production, force the log format to be JSON for structured logging.
	if environment == "production" {
		if l.FORMAT != "json" {
			return fmt.Errorf("%s=production: %s must be \"json\"", envAppEnrionment, envLogFormat)
		}
	}
	return nil
}

// validate ensures logger configuration complies with environment-specific business rules.
func (db *Database) validate(environment string) error {
	errs := []error{
		db.validateDBURL(),
		db.validateDBSSLMode(environment),
		positive(envDBMaxConnLifetime, db.MaxConnLifetime),
		positive(envDBMaxConnIdleTime, db.MaxConnIdleTime),
		positive(envDBConnectTimeout, db.ConnectTimeout),
		positive(envDBStartupTimeout, db.StartupTimeout),
		positive(envDBStatementTimeout, db.StatementTimeout),
		positive(envDBHealthCheckPeriod, db.HealthCheckPeriod),
		positive(envDBReadinessTimeout, db.ReadinessTimeout),
	}

	if db.MaxConns < 1 || db.MaxConns > maxPoolConns {
		errs = append(errs, fmt.Errorf("%s: must be between 1 and %d, got %d", envDBMaxConns, maxPoolConns, db.MaxConns))
	}
	if db.MinConns < 0 || db.MinConns > db.MaxConns {
		errs = append(errs, fmt.Errorf("%s: must be between 0 and %s (%d), got %d",
			envDBMinConns, envDBMaxConns, db.MaxConns, db.MinConns))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (db *Database) validateDBURL() error {
	u, err := url.Parse(db.URL)
	if err != nil {
		return fmt.Errorf("%s: not a valid URL! Provide a valid $s", envDBURL)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("%s: scheme must be postgres:// or postgresql://, got %q", envDBURL, u.Scheme)
	}

	if strings.Trim(u.Path, "/") == "" {
		return fmt.Errorf("%s: database name is missing", envDBURL)
	}

	return nil
}

// validateDBSSLMode ensures a secure database connection mode is configured,
// enforcing strict requirements in production.
func (db *Database) validateDBSSLMode(environment string) error {
	u, err := url.Parse(db.URL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	sslmode := u.Query().Get("sslmode")

	// Enforce strict SSL requirements in production environments
	if environment == "production" {
		switch sslmode {
		case "require", "verify-ca", "verify-full":
			// Valid secure mode
		default:
			if sslmode == "" {
				sslmode = "(not set)"
			}
			return fmt.Errorf("%s: sslmode must be require, verify-ca, or verify-full when %s=production, got %s",
				envDBURL, envAppEnrionment, sslmode)
		}
	}

	return nil
}

func positive(key string, d time.Duration) error {
	if d <= 0 {
		return fmt.Errorf("%s: must be greater than 0, got %v", key, d)
	}
	return nil
}
