package config

import "github.com/kelseyhightower/envconfig"

// Config holds the configuration for the service.
// Values are populated from environment variables.
type Config struct {
	ProjectID string `envconfig:"VACANCIES_GCP_PROJECT"`
	BotToken  string `envconfig:"VACANCIES_BOT_TOKEN"`
	ChatID    string `envconfig:"VACANCIES_CHAT_ID"`
}

// FromEnv loads the configuration from environment variables.
func FromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
