package oauth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestOAuthService_Callback(t *testing.T) {
	var (
		expectedToken        string = faker.Word()
		expectedRefreshToken string = faker.Word()
		// expectedEmail        string = faker.Word()
	)

	type args struct {
		ctx   context.Context
		input oauth.GoogleOauthCallbackInput
	}

	type expected = testutil.Result[oauth.GoogleOauthCallbackOutput]

	type fields struct {
		OauthProvider   oauth.OAuthProviderInterface
		OauthRepoWriter oauth.OAuthRepositoryWriterInterface
		AuthService     auth.AuthServiceInterface
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
				ctx: context.Background(),
				input: oauth.GoogleOauthCallbackInput{
					Code:  faker.Word(),
					State: `actions=["REGISTER_USER"]`,
				},
			},
			expected: expected{
				Err:   nil,
				Value: oauth.GoogleOauthCallbackOutput{Token: expectedToken, RefreshToken: expectedRefreshToken},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				oauthProvider := oauth.NewMockOAuthProviderInterface(ctrl)
				oauthProvider.EXPECT().
					Callback(gomock.Any(), gomock.Any()).
					Return(oauth.GoogleOauthCallbackOutput{Token: expectedToken, RefreshToken: expectedRefreshToken}, nil)

				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					FindAuthWithEmail(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthWithEmailOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: "081234567890",
						Email:   faker.Word(),
					}, nil)
				authService.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Return(auth.UserTokenClaimOutput{
						TokenType: auth.TokenTypeJwt,
						Token:     expectedToken,
					}, nil)
				authService.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					Return(auth.StoreAccessTokenOutput{
						Token:  expectedToken,
						UserId: faker.UUIDHyphenated(),
					}, nil)

				oauthRepoWriter := oauth.NewMockOAuthRepositoryWriterInterface(ctrl)
				oauthRepoWriter.EXPECT().
					CreateAccessToken(gomock.Any(), gomock.Any()).
					Return(oauth.CreateAccessTokenOutput{
						Id: faker.UUIDHyphenated(),
					}, nil)
				oauthRepoWriter.EXPECT().
					BypassCreateNewUser(gomock.Any(), gomock.Any()).
					Return(oauth.BypassCreateNewUserOutput{Id: faker.UUIDHyphenated()}, nil)

				tt.fields.OauthRepoWriter = oauthRepoWriter
				tt.fields.OauthProvider = oauthProvider
				tt.fields.AuthService = authService
			},
		},
		{
			name: "failed generate token",
			args: args{
				ctx:   context.Background(),
				input: oauth.GoogleOauthCallbackInput{},
			},
			expected: expected{
				Err:   oauth.ErrFailedGetCallbackToken,
				Value: oauth.GoogleOauthCallbackOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				oauthProvider := oauth.NewMockOAuthProviderInterface(ctrl)
				oauthProvider.EXPECT().
					Callback(gomock.Any(), gomock.Any()).
					Return(oauth.GoogleOauthCallbackOutput{}, oauth.ErrFailedGetCallbackToken)
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
				Provider:              tt.fields.OauthProvider,
				AuthService:           tt.fields.AuthService,
				OauthRepositoryWriter: tt.fields.OauthRepoWriter,
				Err:                   ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})

			got, err := oauthService.Callback(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
