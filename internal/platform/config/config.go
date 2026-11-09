package config

import "log/slog"

// Environment variable keys for application settings.
const envAppEnrionment = "APP_ENV"

// Environment variable keys for http server settings.
const envHTTPAddr = "HTTP_ADDR"

// Environment variable keys for logger config settings.
const (
	envLogLevel  = "LOG_LEVEL"
	envLogFormat = "LOG_FORMAT"
)

// Config holds all configuration sections for the application.
type Config struct {
	App  *APP
	Http *HTTP
	Logger *Logger
}

// APP holds application-specific settings.
type APP struct {
	Envrionment string
}

// HTTP holds HTTP server settings.
type HTTP struct {
	Addr string
}

// Logger holds logger config settings.
type Logger struct {
	LOG_LEVEL  slog.Level
	LOG_FORMAT string
}
