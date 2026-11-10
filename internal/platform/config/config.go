package config

import (
	"log/slog"
	"time"
)

// Environment variable keys for application settings.
const envAppEnrionment = "APP_ENV"

// Environment variable keys for http server settings.
const (
	envHTTPAddr              = "HTTP_ADDR"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	envHTTPReadTimeout       = "HTTP_READ_TIMEOUT"
	envHTTPWriteTimeout      = "HTTP_WRITE_TIMEOUT"
	envHTTPIdleTimeout       = "HTTP_IDLE_TIMEOUT"
	envHTTPDrainDelay        = "HTTP_DRAIN_DELAY"
	envHTTPShutdownTimeout   = "HTTP_SHUTDOWN_TIMEOUT"
	envHTTPRequestTimeout    = "HTTP_REQUEST_TIMEOUT"
	envHTTPMaxBodyBytes      = "HTTP_MAX_BODY_BYTES"
)

// Environment variable keys for logger config settings.
const (
	envLogLevel  = "LOG_LEVEL"
	envLogFormat = "LOG_FORMAT"
	envAddSource = "LOG_ADD_SOURCE"
)

// Config holds all configuration sections for the application.
type Config struct {
	App    *APP
	Http   *HTTP
	Logger *Logger
}

// APP holds application-specific settings.
type APP struct {
	Envrionment string
}

// HTTP holds HTTP server settings.
type HTTP struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	DrainDelay        time.Duration
	ShutdownTimeout   time.Duration
	RequestTimeout    time.Duration
	MaxBodyBytes      int64
}

// Logger holds logger config settings.
type Logger struct {
	LEVEL     slog.Level
	FORMAT    string
	AddSource bool
}
