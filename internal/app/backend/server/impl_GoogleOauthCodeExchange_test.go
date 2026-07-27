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

func TestServer_GoogleOauthCodeExchange(t *testing.T) {
	var (
		expectedUrl string = faker.Word()
	)

	type args struct {
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
		mock     func(ctx echo.Context, tt *test, mock *gomock.Controller)
	}

	tests := []test{
		{
			name: "success with 200 response",
			args: args{},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				oauthService := oauth.NewMockOAuthServiceInterface(ctrl)
				oauthService.EXPECT().
					ExchangeToken(ctx.Request().Context(), gomock.Any()).
					Return(oauth.GoogleOauthCodeExchangeOutput{Url: expectedUrl}, nil)

				tt.fields.OAuthService = oauthService

				result, _ := json.Marshal(api.GoogleOauthCodeExchangeResponse{
					Url: expectedUrl,
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
			svc := server.NewServer(server.ServerOpts{
				OAuthService: tt.fields.OAuthService,
			})

			svc.GoogleOauthCodeExchange(ctx)
			testutil.RequireHttpResultJson(t, tt.expected, rec)
		})
	}
}
