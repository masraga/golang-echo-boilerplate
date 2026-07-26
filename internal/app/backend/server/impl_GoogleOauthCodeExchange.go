package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

// Oauth code exchange
// (POST /api/v1/oauth)
func (s *Server) GoogleOauthCodeExchange(ctx echo.Context) error {
	output, err := s.OAuthService.ExchangeToken(ctx.Request().Context(), oauth.GoogleOauthCodeExchangeInput{})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.GoogleOauthCodeExchangeResponse{
		Url: output.Url,
	})
}
