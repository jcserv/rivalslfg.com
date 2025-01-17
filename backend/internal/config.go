package internal

import (
	"errors"

	"github.com/jcserv/rivalslfg/internal/utils/env"
)

type Configuration struct {
	Region       string
	Environment  string
	HTTPPort     string
	WSPort       string
	DatabaseURL  string
	CacheURL     string
	JWTSecretKey string
}

func NewConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	cfg.Region = env.GetString("REGION", "us-east-1")
	cfg.Environment = env.GetString("ENVIRONMENT", "dev")
	cfg.HTTPPort = env.GetString("HTTP_PORT", "8080")
	cfg.WSPort = env.GetString("WS_PORT", "8081")
	cfg.DatabaseURL = env.GetString("DATABASE_URL", "")
	cfg.CacheURL = env.GetString("CACHE_URL", "")
	cfg.JWTSecretKey = env.GetString("JWT_SECRET_KEY", "")
	return cfg, nil
}

func (c *Configuration) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.CacheURL == "" {
		return errors.New("CACHE_URL is required")
	}
	if c.JWTSecretKey == "" {
		return errors.New("JWT_SECRET_KEY is required")
	}
	return nil
}
