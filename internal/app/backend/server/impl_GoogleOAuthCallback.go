package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

// Google OAuth callback
// (GET /api/v1/oauth/callback)
func (s *Server) GoogleOAuthCallback(ctx echo.Context, params api.GoogleOAuthCallbackParams) error {
	output, err := s.OAuthService.Callback(ctx.Request().Context(), oauth.GoogleOauthCallbackInput{
		Code:  params.Code,
		State: params.State,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	return returnOk(ctx, api.GoogleOAuthCallbackResponse{
		Token:        output.Token,
		RefreshToken: output.RefreshToken,
		State:        &output.State,
	})
}
