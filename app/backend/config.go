package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/masraga/golang-echo-boilerplate/external/fcm"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

type Config struct {
	ShowErrMode               ctxerr.ShowErrMode         `env:"SHOW_ERR_MODE"`
	AppPort                   int64                      `env:"APP_PORT"`
	DatabaseUrl               string                     `env:"DATABASE_URL"`
	JwtSecret                 auth.JwtSecretType         `env:"JWT_SECRET"`
	JwtExpiration             auth.JwtExpirationType     `env:"JWT_EXPIRATION"`
	CryptoKey                 crypto.ConfigCryptoKey     `env:"CRYPTO_KEY"`
	AuthAccessBootstrapUserId string                     `env:"AUTH_ACCESS_BOOTSTRAP_USER_ID"`
	FcmServiceAccountId       fcm.ConfigServiceAccountId `env:"FCM_SERVICE_ACCOUNT_ID"`
}

func LoadConfig() *Config {
	var cfg Config

	_ = cleanenv.ReadConfig(".env", &cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil
	}

	return &cfg
}
