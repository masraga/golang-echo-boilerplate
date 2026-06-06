package main

import (
	"context"

	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func Initialize(ctx context.Context) (cfg *Config, server *server.Server) {
	cfg = LoadConfig()
	sv, err := InitializeService(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if cfg.AuthAccessBootstrapUserId != "" {
		_, err = sv.AuthService.BootstrapUserApiContracts(context.Background(), auth.BootstrapUserApiContractsInput{
			UserId: cfg.AuthAccessBootstrapUserId,
		})
		if err != nil {
			panic(err)
		}
	}
	return cfg, sv
}
