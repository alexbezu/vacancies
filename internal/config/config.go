package config

import "github.com/kelseyhightower/envconfig"

// Config holds the configuration for the service.
// Values are populated from environment variables.
type Config struct {
	Port     string   `envconfig:"PORT" default:"8080"`
	WebHook  string   `envconfig:"WEBHOOK"`
	Firebase struct {
		ProjectID string `envconfig:"FIREBASE_PROJECT_ID"`
	}
}

// FromEnv loads the configuration from environment variables.
func FromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}