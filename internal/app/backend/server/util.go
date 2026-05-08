package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) returnOk(e echo.Context, data interface{}) error {
	return e.JSON(200, data)
}

func (s *Server) bindOrReturnBadRequest(e echo.Context, i any) error {
	if err := e.Bind(&i); err != nil {
		return e.JSON(http.StatusBadRequest, "")
	}
	return nil
}
