package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
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

			authToken := strings.TrimPrefix(token, "Bearer ")
			if authToken == token || authToken == "" {
				return returnUnauthorized(c)
			}

			authOutput, err := authInterface.ValidateJwtToken(c.Request().Context(), auth.ValidateJwtTokenInput{
				Token: authToken,
			})
			if err != nil {
				return returnUnauthorized(c)
			}

			_, err = authInterface.ValidateUserApiContract(c.Request().Context(), auth.ValidateUserApiContractInput{
				UserId:         authOutput.UserId,
				EndpointPath:   endpointPath,
				EndpointMethod: method,
			})
			if err != nil {
				return returnForbidden(c)
			}

			return next(c)
		}
	}
}
