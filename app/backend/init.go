package main

import (
	"github.com/masraga/kerp-api/internal/app/backend/server"
)

func Initialize() (cfg *Config, server *server.Server) {
	cfg = LoadConfig()
	sv, err := InitializeService(cfg)
	if err != nil {
		panic(err)
	}
	return cfg, sv
}
