package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func returnUnauthorized(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
		"error": "Unauthorized",
	})
}

func returnForbidden(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, map[string]interface{}{
		"error": "Forbidden",
	})
}

func skipValidation(path string, method string) bool {
	if filter, ok := skipAuthValidationFilterMap[path]; ok {
		for _, m := range filter.Method {
			if m == method {
				return true
			}
		}
	}
	return false
}
