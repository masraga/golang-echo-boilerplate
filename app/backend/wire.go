//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/service/notification"
)

func InitializeService(ctx context.Context, config *Config) (*server.Server, error) {
	wire.Build(
		ProvideSqlDb,
		ProvideDbTx,
		ProvideSqlDialect,
		ProvideAuthAccessBootstrapUserId,
		ProvideZerolog,
		ProvidePushNotificationService,

		wire.FieldsOf(new(*Config),
			"JwtSecret",
			"JwtExpiration",
			"ShowErrMode",
			"CryptoKey",
			"AuthAccessBootstrapUserId",
		),

		wire.Struct(new(ctxerr.CtxErrOpts), "*"),
		wire.Struct(new(crypto.CryptoServiceOpts), "*"),
		crypto.NewCryptoService,
		wire.Bind(new(crypto.CryptoServiceInterface), new(*crypto.CryptoService)),

		// authentication
		wire.Struct(new(auth.AuthRepositoryOpts), "*"),
		auth.NewAuthRepository,
		wire.Bind(new(auth.AuthRepositoryWriterInterface), new(*auth.AuthRepository)),
		wire.Bind(new(auth.AuthRepositoryReaderInterface), new(*auth.AuthRepository)),
		ctxerr.NewCtxErr,
		wire.Struct(new(auth.AuthServiceOpts), "*"),
		auth.NewAuthService,
		wire.Bind(new(auth.AuthServiceInterface), new(*auth.AuthService)),

		// notification
		wire.Struct(new(notification.NotificationServiceOpts), "*"),
		notification.NewNotificationService,
		wire.Bind(new(notification.NotificationServiceInterface), new(*notification.NotificationService)),

		wire.Struct(new(server.ServerOpts), "*"),
		server.NewServer,
	)

	return &server.Server{}, nil
}
