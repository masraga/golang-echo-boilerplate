package oauth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestOAuthService_ExchangeToken(t *testing.T) {
	var (
		expectedUrl string = faker.Word()
	)

	type args struct {
		ctx   context.Context
		input oauth.GoogleOauthCodeExchangeInput
	}

	type expected = testutil.Result[oauth.GoogleOauthCodeExchangeOutput]

	type fields struct {
		OauthProvider oauth.OAuthProviderInterface
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
			name: "success generate token",
			args: args{
				ctx:   context.Background(),
				input: oauth.GoogleOauthCodeExchangeInput{},
			},
			expected: expected{
				Err:   nil,
				Value: oauth.GoogleOauthCodeExchangeOutput{Url: expectedUrl},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				oauthProvider := oauth.NewMockOAuthProviderInterface(ctrl)
				oauthProvider.EXPECT().
					ExchangeToken(gomock.Any(), gomock.Any()).
					Return(oauth.GoogleOauthCodeExchangeOutput{Url: expectedUrl}, nil)
				tt.fields.OauthProvider = oauthProvider
			},
		},
		{
			name: "failed generate token",
			args: args{
				ctx:   context.Background(),
				input: oauth.GoogleOauthCodeExchangeInput{},
			},
			expected: expected{
				Err:   oauth.ErrFailedExchangeToken,
				Value: oauth.GoogleOauthCodeExchangeOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				oauthProvider := oauth.NewMockOAuthProviderInterface(ctrl)
				oauthProvider.EXPECT().
					ExchangeToken(gomock.Any(), gomock.Any()).
					Return(oauth.GoogleOauthCodeExchangeOutput{}, oauth.ErrFailedExchangeToken)
				tt.fields.OauthProvider = oauthProvider
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

			oauthService := oauth.NewOAuthService(oauth.OAuthServiceOpts{
				Provider: tt.fields.OauthProvider,
				Err:      ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})

			got, err := oauthService.ExchangeToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
