package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/internal/util/traceerr"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	req := c.Request()
	c.Logger().Errorf(
		"request error: method=%s path=%s route=%s request_id=%s error=%v origin=%s\nstacktrace:\n%s",
		req.Method,
		req.URL.RequestURI(),
		c.Path(),
		requestID(c),
		err,
		errorOrigin(err),
		debug.Stack(),
	)

	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func requestID(c echo.Context) string {
	requestID := c.Request().Header.Get(echo.HeaderXRequestID)
	if requestID != "" {
		return requestID
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func errorOrigin(err error) string {
	file, line, ok := traceerr.Location(err)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
