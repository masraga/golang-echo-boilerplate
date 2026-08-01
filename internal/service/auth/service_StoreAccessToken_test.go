package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_StoreAccessToken(t *testing.T) {
	var (
		expectedToken     string = faker.Word()
		expectedUserId    string = faker.UUIDHyphenated()
		expectedExpiredAt int64  = time.Now().Add(time.Hour).UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input auth.StoreAccessTokenInput
	}

	type expected = testutil.Result[auth.StoreAccessTokenOutput]

	type fields struct {
		AuthRepositoryWriter auth.AuthRepositoryWriterInterface
	}

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success store access token",
			args: args{
				ctx: context.Background(),
				input: auth.StoreAccessTokenInput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.StoreAccessTokenOutput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
					IsActive:      true,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					Return(auth.StoreAccessTokenOutput{
						Token:         expectedToken,
						UserId:        expectedUserId,
						ExpiredAtUtc0: expectedExpiredAt,
						IsActive:      true,
					}, nil)

				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})
			got, err := authService.StoreAccessToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
