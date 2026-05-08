package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/masraga/kerp-api/internal/service/auth"
)

type Config struct {
	AppPort       int64                  `env:"APP_PORT"`
	DatabaseUrl   string                 `env:"DATABASE_URL"`
	JwtSecret     auth.JwtSecretType     `env:"JWT_SECRET"`
	JwtExpiration auth.JwtExpirationType `env:"JWT_EXPIRATION"`
}

func LoadConfig() *Config {
	var cfg Config

	_ = cleanenv.ReadConfig(".env", &cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil
	}

	return &cfg
}
