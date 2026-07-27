package oauth

import (
	"context"
)

type OAuthServiceInterface interface {
	ExchangeToken(ctx context.Context, input GoogleOauthCodeExchangeInput) (output GoogleOauthCodeExchangeOutput, err error)
	Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error)
}

type OAuthProviderInterface interface {
	ExchangeToken(ctx context.Context, input GoogleOauthCodeExchangeInput) (output GoogleOauthCodeExchangeOutput, err error)
	Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error)
}

type OAuthRepositoryReaderInterface interface {
	GetAccessTokenById(ctx context.Context, input GetAccessTokenByIdInput) (output GetAccessTokenByIdOutput, err error)
}

type OAuthRepositoryWriterInterface interface{}
