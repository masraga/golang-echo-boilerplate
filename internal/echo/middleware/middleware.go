package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func AuthValidationFilterMiddleware(auth auth.AuthServiceInterface) echo.MiddlewareFunc {
	return authValidationFilter(auth)
}
