package auth

import (
	"context"

	"github.com/masraga/kerp-api/internal/dbtx"
)

type AuthServiceInterface interface {
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
	CreateToken(ctx context.Context, input UserTokenClaimInput) (output UserTokenClaimOutput, err error)
	CreateJWTToken(ctx context.Context, input CreateJWTTokenInput) (output CreateJWTTokenOutput, err error)
}

type AuthRepositoryWriterInterface interface {
	dbtx.DbTxInterface
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
}

type AuthRepositoryReaderInterface interface {
	FindAuth(ctx context.Context, input FindAuthInput) (output FindAuthOutput, err error)
}
