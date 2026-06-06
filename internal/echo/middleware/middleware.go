package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func AuthValidationFilterMiddleware(auth auth.AuthServiceInterface) echo.MiddlewareFunc {
	return authValidationFilter(auth)
}
