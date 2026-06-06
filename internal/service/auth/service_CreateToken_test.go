package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.UserTokenClaimInput
	}

	type expected = testutil.Result[auth.UserTokenClaimOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(t *testing.T, gomock *gomock.Controller)
	}

	tests := []test{
		{
			name: "should fail when use jwt provider",
			args: args{
				ctx: context.Background(),
				input: auth.UserTokenClaimInput{
					TokenType: auth.TokenTypeJwt,
					UserId:    "",
				},
			},
			expected: expected{
				Err:   auth.ErrClaimJwtToken,
				Value: auth.UserTokenClaimOutput{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewAuthService(
				auth.AuthServiceOpts{
					JwtSecret:     auth.JwtSecretType(faker.Word()),
					JwtExpiration: auth.JwtExpirationType(60),
					Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				},
			)

			res, err := svc.CreateToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
