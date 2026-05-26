package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateJWTTOken(t *testing.T) {

	type args struct {
		ctx   context.Context
		input auth.CreateJWTTokenInput
	}

	type expected = testutil.Result[auth.CreateJWTTokenOutput]

	type test struct {
		name     string
		expected expected
		args     args
		mock     func(t *testing.T, mock *gomock.Controller)
	}

	tests := []test{
		{
			name: "failed when userId not provided",
			expected: expected{
				Err:   auth.ErrClaimJwtToken,
				Value: auth.CreateJWTTokenOutput{},
			},
			args: args{
				ctx: context.Background(),
				input: auth.CreateJWTTokenInput{
					UserId:        faker.UUIDHyphenated(),
					ExpiredAtUtc0: 0,
				},
			},
		},
		{
			name: "failed when expiredAtUtc0 not provided",
			expected: expected{
				Err:   auth.ErrClaimJwtToken,
				Value: auth.CreateJWTTokenOutput{},
			},
			args: args{
				ctx: context.Background(),
				input: auth.CreateJWTTokenInput{
					UserId: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(t, ctrl)
			}

			svc := auth.NewAuthService(
				auth.AuthServiceOpts{
					Err: ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				},
			)

			res, err := svc.CreateJWTToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
