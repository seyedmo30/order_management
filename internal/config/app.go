/*
Package config contains the configuration for the application.
*/
package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

// App holds the application configuration.
type App struct {
	DatabaseConfig
	ServiceConfig
}

// DatabaseConfig contains the database connection details.
type DatabaseConfig struct {
	LogLevel string `env:"DB_LOG_LEVEL" envDefault:"INFO"`
}

// ServiceConfig contains application-specific settings.
type ServiceConfig struct {
	ReportInterval      int `env:"SERVICE_REPORT_INTERVAL" envDefault:"2"`
	WorkerCount         int `env:"SERVICE_WORKER_COUNT" envDefault:"5"`
	OrderProcessTimeout int `env:"SERVICE_ORDER_PROCESS_TIMEOUT" envDefault:"5"`
}

// Load initializes and loads the configuration from environment variables.
func Load() App {
	var cfg App

	// Load environment variables into the App struct
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	return cfg
}
