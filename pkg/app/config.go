package app

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port        string `env:"PORT"`
	DatabaseURL string `env:"DATABASE_URL, required"`
	RiotApiKey  string `env:"RIOT_API_KEY"`
}

func NewConfig(ctx context.Context, cfg *Config) error {
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	return nil
}
