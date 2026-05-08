package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/util/pointer"
)

func (s *Server) GetPing(ctx echo.Context) error {
	return s.returnOk(ctx, api.GetPingResponse{
		Message: pointer.String("PONG"),
	})
}
