package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

// Oauth code exchange
// (POST /api/v1/oauth/google)
func (s *Server) GoogleOauthCodeExchange(ctx echo.Context) error {
	var req api.GoogleOauthCodeExchangeRequest
	if err := bindOrReturnBadRequest(ctx, &req); err != nil {
		return err
	}
	var actions []oauth.ExchangeTokenAction
	for _, act := range req.Actions {
		actions = append(actions, oauth.ExchangeTokenAction(act))
	}
	output, err := s.OAuthService.ExchangeToken(ctx.Request().Context(), oauth.GoogleOauthCodeExchangeInput{
		Actions: actions,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.GoogleOauthCodeExchangeResponse{
		Url: output.Url,
	})
}
