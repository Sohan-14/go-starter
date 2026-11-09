package config

// Environment variable keys for application settings.
const envAppEnrionment = "APP_ENV"

// Environment variable keys for http server settings.
const envHTTPAddr = "HTTP_ADDR"

// Config holds all configuration sections for the application.
type Config struct {
	App  *APP
	Http *HTTP
}

// APP holds application-specific settings.
type APP struct {
	Envrionment string
}

// HTTP holds HTTP server settings.
type HTTP struct {
	Addr string
}
