//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/masraga/golang-echo-boilerplate/external/gauth"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/service/notification"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
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

		// google oauth
		wire.FieldsOf(new(*Config),
			"GoogleAuthClientId",
			"GoogleAuthClientSecret",
			"GoogleAuthCallbackUrl",
		),
		// wire.Struct(new(oauth.OAuthRepositoryOpts), "*"),
		// oauth.NewOAuthRepository,
		// wire.Bind(new(oauth.OAuthRepositoryReaderInterface), new(*oauth.OAuthRepository)),
		// wire.Bind(new(oauth.OAuthRepositoryWriterInterface), new(*oauth.OAuthRepository)),
		wire.Struct(new(gauth.GAuthServiceOpts), "*"),
		gauth.NewGAuthService,
		wire.Bind(new(oauth.OAuthProviderInterface), new(*gauth.GAuthService)),
		wire.Struct(new(oauth.OAuthServiceOpts), "*"),
		oauth.NewOAuthService,
		wire.Bind(new(oauth.OAuthServiceInterface), new(*oauth.OAuthService)),

		wire.Struct(new(server.ServerOpts), "*"),
		server.NewServer,
	)

	return &server.Server{}, nil
}
