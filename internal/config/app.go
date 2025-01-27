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
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

// ServiceConfig contains application-specific settings.
type ServiceConfig struct {
	Port string `env:"SERVICE_PORT" envDefault:"8089"`
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
