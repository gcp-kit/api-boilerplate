package config

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/heetch/confita"
)

type Config struct {
}

var defaultConfig = Config{}

// ReadConfig - read config from environment variables
func ReadConfig(ctx context.Context) (*Config, error) {
	loader := confita.NewLoader()

	cfg := defaultConfig
	if err := loader.Load(ctx, &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	return &cfg, nil
}
