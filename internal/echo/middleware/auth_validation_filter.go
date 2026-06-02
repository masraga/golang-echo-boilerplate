package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func authValidationFilter(authInterface auth.AuthServiceInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			endpointPath := c.Path()
			method := c.Request().Method
			token := c.Request().Header.Get("Authorization")

			if skipValidation(endpointPath, method) {
				return next(c)
			}

			if token == "" {
				return returnUnauthorized(c)
			}

			authToken := strings.Split(token, "Bearer ")[1]

			if _, err := authInterface.ValidateJwtToken(c.Request().Context(), auth.ValidateJwtTokenInput{
				Token: authToken,
			}); err != nil {
				return returnUnauthorized(c)
			}

			return next(c)
		}
	}
}
