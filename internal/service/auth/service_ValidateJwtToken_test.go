package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_ValidateJwtToken(t *testing.T) {
	var (
		jwtSecret         = auth.JwtSecretType(faker.Word())
		expectedUserId    = faker.UUIDHyphenated()
		expectedExpiredAt = time.Now().Add(time.Hour).UnixMilli()
		expectedIssuedAt  = time.Now().UnixMilli()
		validToken        = createJwtToken(t, jwtSecret, expectedUserId, expectedExpiredAt, expectedIssuedAt)
		expiredAt         = time.Now().Add(-time.Hour).UnixMilli()
		expiredToken      = createJwtToken(t, jwtSecret, expectedUserId, expiredAt, expectedIssuedAt)
		invalidSignToken  = createJwtToken(t, auth.JwtSecretType(faker.Word()), expectedUserId, expectedExpiredAt, expectedIssuedAt)
	)

	type args struct {
		ctx   context.Context
		input auth.ValidateJwtTokenInput
	}

	type fields struct {
		AuthRepositoryReader *auth.MockAuthRepositoryReaderInterface
	}

	type expected = testutil.Result[auth.ValidateJwtTokenOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should success validate jwt token with active stored access token",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateJwtTokenInput{
					Token: validToken,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.ValidateJwtTokenOutput{
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
					IssuerAtUtc0:  expectedIssuedAt,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAccessToken(gomock.Any(), auth.FindAccessTokenInput{
						Token:  validToken,
						UserId: expectedUserId,
					}).
					Return(auth.FindAccessTokenOutput{
						Token:         validToken,
						UserId:        expectedUserId,
						ExpiredAtUtc0: expectedExpiredAt,
						IsActive:      true,
					}, nil)
				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed when active stored access token not found",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateJwtTokenInput{
					Token: validToken,
				},
			},
			expected: expected{
				Err:   auth.ErrFindAccessTokenNotFound,
				Value: auth.ValidateJwtTokenOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAccessToken(gomock.Any(), auth.FindAccessTokenInput{
						Token:  validToken,
						UserId: expectedUserId,
					}).
					Return(auth.FindAccessTokenOutput{}, auth.ErrFindAccessTokenNotFound)
				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed when jwt token expired",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateJwtTokenInput{
					Token: expiredToken,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthTokenExpired,
				Value: auth.ValidateJwtTokenOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				tt.fields.AuthRepositoryReader = auth.NewMockAuthRepositoryReaderInterface(ctrl)
			},
		},
		{
			name: "should failed when jwt token signature invalid",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateJwtTokenInput{
					Token: invalidSignToken,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthSigInvalid,
				Value: auth.ValidateJwtTokenOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				tt.fields.AuthRepositoryReader = auth.NewMockAuthRepositoryReaderInterface(ctrl)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            jwtSecret,
				JwtExpiration:        auth.JwtExpirationType(60),
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
			})

			got, err := authService.ValidateJwtToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}

func createJwtToken(t *testing.T, jwtSecret auth.JwtSecretType, userId string, expiredAtUtc0 int64, issuerAtUtc0 int64) string {
	t.Helper()

	authService := auth.NewAuthService(auth.AuthServiceOpts{
		JwtSecret:     jwtSecret,
		JwtExpiration: auth.JwtExpirationType(60),
		Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
	})

	output, err := authService.CreateJWTToken(context.Background(), auth.CreateJWTTokenInput{
		ExpiredAtUtc0: expiredAtUtc0,
		IssuerAtUtc0:  issuerAtUtc0,
		UserId:        userId,
	})
	require.NoError(t, err)

	return output.Token
}
