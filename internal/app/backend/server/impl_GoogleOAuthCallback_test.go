package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_GoogleOAuthCallback(t *testing.T) {
	var (
		expectedCode         string = faker.Word()
		expectedToken        string = faker.Word()
		expectedRefreshToken string = faker.Word()
	)

	type args struct {
		input api.GoogleOAuthCallbackParams
	}

	type fields struct {
		OAuthService oauth.OAuthServiceInterface
	}

	type expected = testutil.HttpResult

	type test struct {
		name     string
		expected expected
		fields   fields
		args     args
		mock     func(ctx echo.Context, tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success with 200",
			args: args{
				input: api.GoogleOAuthCallbackParams{Code: expectedCode},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				oauthService := oauth.NewMockOAuthServiceInterface(ctrl)
				oauthService.EXPECT().
					Callback(ctx.Request().Context(), gomock.Any()).
					Return(oauth.GoogleOauthCallbackOutput{
						Token:        expectedToken,
						RefreshToken: expectedRefreshToken,
					}, nil)

				tt.fields.OAuthService = oauthService

				result, _ := json.Marshal(api.GoogleOAuthCallbackResponse{
					RefreshToken: expectedRefreshToken,
					Token:        expectedToken,
				})
				tt.expected.Code = http.StatusOK
				tt.expected.Body = string(result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			if tt.mock != nil {
				tt.mock(ctx, &tt, ctrl)
			}
			server := server.NewServer(server.ServerOpts{
				OAuthService: tt.fields.OAuthService,
			})

			server.GoogleOAuthCallback(ctx, tt.args.input)
			testutil.RequireHttpResultJson(t, tt.expected, rec)
		})
	}
}
