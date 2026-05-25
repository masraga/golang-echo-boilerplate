//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func InitializeService(config *Config) (*server.Server, error) {
	wire.Build(
		ProvideSqlDb,
		ProvideDbTx,
		ProvideSqlDialect,
		ProvideZerolog,
		wire.Struct(new(ctxerr.CtxErrOpts), "*"),

		wire.FieldsOf(new(*Config), "JwtSecret", "JwtExpiration"),
		wire.Struct(new(auth.AuthRepositoryOpts), "*"),
		auth.NewAuthRepository,
		wire.Bind(new(auth.AuthRepositoryWriterInterface), new(*auth.AuthRepository)),
		wire.Bind(new(auth.AuthRepositoryReaderInterface), new(*auth.AuthRepository)),
		ctxerr.NewCtxErr,
		wire.Struct(new(auth.AuthServiceOpts), "*"),
		auth.NewAuthService,
		wire.Bind(new(auth.AuthServiceInterface), new(*auth.AuthService)),
		server.NewServer,
	)

	return &server.Server{}, nil
}
