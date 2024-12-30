package internal

import (
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
	cfg.Environment = env.GetString("ENVIRONMENT", "prod")
	cfg.HTTPPort = env.GetString("HTTP_PORT", "8080")
	return cfg, nil
}
