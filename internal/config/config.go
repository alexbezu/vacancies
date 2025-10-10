package config

import "github.com/kelseyhightower/envconfig"

// Config holds the configuration for the service.
// Values are populated from environment variables.
type Config struct {
	WebHook   string `envconfig:"WEBHOOK"`
	ProjectID string `envconfig:"VACANCIES_GCP_PROJECT"`
}

// FromEnv loads the configuration from environment variables.
func FromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
