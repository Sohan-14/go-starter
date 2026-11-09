package config

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

// validate runs business logic validation across all configuration sections.
func (c *Config) validate() error {
	var errs []error

	if err := c.App.validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Http.validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Logger.validate(c.App.Envrionment); err != nil {
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
		if l.LOG_FORMAT != "json" {
			return fmt.Errorf("%s=production: %s must be \"json\"", envAppEnrionment, envLogFormat)
		}
	}
	return nil
}
