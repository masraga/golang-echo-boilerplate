package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
)

func (s *Server) GetPing(ctx echo.Context) error {
	return returnOk(ctx, api.GetPingResponse{
		Message: pointer.String("PONG"),
	})
}
