package internal

import (
	"errors"

	"github.com/jcserv/rivalslfg/internal/utils/env"
)

type Configuration struct {
	Region      string
	Environment string
	HTTPPort    string
	DatabaseURL string
}

func NewConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	cfg.Region = env.GetString("REGION", "us-east-1")
	cfg.Environment = env.GetString("ENVIRONMENT", "dev")
	cfg.HTTPPort = env.GetString("HTTP_PORT", "8080")
	cfg.DatabaseURL = env.GetString("DATABASE_URL", "")
	return cfg, nil
}

func (c *Configuration) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	return nil
}