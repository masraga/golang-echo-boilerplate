package oauth_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestOAuthService_GetAccessTokenById(t *testing.T) {
	var (
		expectedId                 string = faker.UUIDHyphenated()
		expectedUsername           string = faker.Email()
		expectedEmail              string = faker.Email()
		expectedAccessToken        string = faker.Word()
		expectedAccessRefreshToken string = faker.Word()
		expectedExpiryAtUtc0       int64  = time.Now().UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input oauth.GetAccessTokenByIdInput
	}

	type expected = testutil.Result[oauth.GetAccessTokenByIdOutput]

	type fields struct {
		oauthRepoReader oauth.OAuthRepositoryReaderInterface
	}

	type test struct {
		name     string
		expected expected
		args     args
		mock     func(tt *test, ctrl *gomock.Controller)
		fields   fields
	}

	tests := []test{
		{
			name: "success get access token by id",
			args: args{
				ctx:   context.Background(),
				input: oauth.GetAccessTokenByIdInput{},
			},
			expected: expected{
				Err: nil,
				Value: oauth.GetAccessTokenByIdOutput{
					Id:           expectedId,
					Username:     pointer.String(expectedUsername),
					Email:        pointer.String(expectedEmail),
					AccessToken:  expectedAccessToken,
					RefreshToken: expectedAccessRefreshToken,
					ExpiryAtUtc0: pointer.Int64(expectedExpiryAtUtc0),
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				oauthRepoReader := oauth.NewMockOAuthRepositoryReaderInterface(ctrl)
				oauthRepoReader.EXPECT().
					GetAccessTokenById(gomock.Any(), gomock.Any()).
					Return(oauth.GetAccessTokenByIdOutput{
						Id:           expectedId,
						Username:     pointer.String(expectedUsername),
						Email:        pointer.String(expectedEmail),
						AccessToken:  expectedAccessToken,
						RefreshToken: expectedAccessRefreshToken,
						ExpiryAtUtc0: pointer.Int64(expectedExpiryAtUtc0),
					}, nil)
				tt.fields.oauthRepoReader = oauthRepoReader
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

			oauthService := oauth.NewOAuthService(oauth.OAuthServiceOpts{
				OauthRepositoryReader: tt.fields.oauthRepoReader,
				Err:                   ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})

			got, err := oauthService.GetAccessTokenById(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
