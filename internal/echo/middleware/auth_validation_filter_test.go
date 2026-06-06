package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthValidationFilterMiddleware(t *testing.T) {
	type args struct {
		path          string
		method        string
		authorization string
	}

	type fields struct {
		authService auth.AuthServiceInterface
	}

	type expected struct {
		status     int
		nextCalled bool
	}

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should skip public route",
			args: args{
				path:   "/api/v1/auth/register/phone",
				method: http.MethodPost,
			},
			expected: expected{
				status:     http.StatusNoContent,
				nextCalled: true,
			},
		},
		{
			name: "should return unauthorized when token empty",
			args: args{
				path:   "/api/v1/auth/api-contracts",
				method: http.MethodGet,
			},
			expected: expected{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "should return unauthorized when bearer token malformed",
			args: args{
				path:          "/api/v1/auth/api-contracts",
				method:        http.MethodGet,
				authorization: "token",
			},
			expected: expected{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "should return unauthorized when jwt invalid",
			args: args{
				path:          "/api/v1/auth/api-contracts",
				method:        http.MethodGet,
				authorization: "Bearer access-token",
			},
			expected: expected{
				status: http.StatusUnauthorized,
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					ValidateJwtToken(gomock.Any(), auth.ValidateJwtTokenInput{Token: "access-token"}).
					Return(auth.ValidateJwtTokenOutput{}, errors.New("invalid token"))
				tt.fields.authService = authService
			},
		},
		{
			name: "should return forbidden when user has no api access",
			args: args{
				path:          "/api/v1/auth/api-contracts",
				method:        http.MethodGet,
				authorization: "Bearer access-token",
			},
			expected: expected{
				status: http.StatusForbidden,
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					ValidateJwtToken(gomock.Any(), auth.ValidateJwtTokenInput{Token: "access-token"}).
					Return(auth.ValidateJwtTokenOutput{UserId: "user-id"}, nil)

				authService.EXPECT().
					ValidateUserApiContract(gomock.Any(), auth.ValidateUserApiContractInput{
						UserId:         "user-id",
						EndpointPath:   "/api/v1/auth/api-contracts",
						EndpointMethod: http.MethodGet,
					}).
					Return(auth.ValidateUserApiContractOutput{}, auth.ErrUserApiContractForbidden)

				tt.fields.authService = authService
			},
		},
		{
			name: "should call next when jwt and api access valid",
			args: args{
				path:          "/api/v1/auth/api-contracts",
				method:        http.MethodGet,
				authorization: "Bearer access-token",
			},
			expected: expected{
				status:     http.StatusNoContent,
				nextCalled: true,
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					ValidateJwtToken(gomock.Any(), auth.ValidateJwtTokenInput{Token: "access-token"}).
					Return(auth.ValidateJwtTokenOutput{UserId: "user-id"}, nil)

				authService.EXPECT().
					ValidateUserApiContract(gomock.Any(), auth.ValidateUserApiContractInput{
						UserId:         "user-id",
						EndpointPath:   "/api/v1/auth/api-contracts",
						EndpointMethod: http.MethodGet,
					}).
					Return(auth.ValidateUserApiContractOutput{IsAllowed: true}, nil)

				tt.fields.authService = authService
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

			e := echo.New()
			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			if tt.args.authorization != "" {
				req.Header.Set("Authorization", tt.args.authorization)
			}
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath(tt.args.path)

			nextCalled := false
			next := func(c echo.Context) error {
				nextCalled = true
				return c.NoContent(http.StatusNoContent)
			}

			err := AuthValidationFilterMiddleware(tt.fields.authService)(next)(ctx)
			require.NoError(t, err)
			require.Equal(t, tt.expected.status, rec.Code)
			require.Equal(t, tt.expected.nextCalled, nextCalled)
		})
	}
}
