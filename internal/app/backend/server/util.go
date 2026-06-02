package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
)

func returnBadRequest(e echo.Context, data interface{}) error {
	return e.JSON(http.StatusBadRequest, api.ErrorBadRequest{Error: data.(string)})
}

func returnCreated(e echo.Context, data interface{}) error {
	return e.JSON(http.StatusCreated, data)
}

func returnOk(e echo.Context, data interface{}) error {
	return e.JSON(http.StatusOK, data)
}

func returnNotImplemented(e echo.Context) error {
	return e.JSON(http.StatusNotImplemented, "not implemented")
}

func bindOrReturnBadRequest(e echo.Context, i any) error {
	if err := e.Bind(i); err != nil {
		return e.JSON(http.StatusBadRequest, "")
	}
	return nil
}

func returnError(e echo.Context, err error) error {
	status := http.StatusInternalServerError
	for target, mappedStatus := range mapError {
		if errors.Is(err, target) {
			status = mappedStatus
			break
		}
	}
	if status == http.StatusInternalServerError {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "unknown error"})
	}

	switch status {
	case http.StatusBadRequest:
		return e.JSON(status, api.ErrorBadRequest{Error: err.Error()})
	case http.StatusUnauthorized:
		return e.JSON(status, api.ErrorUnauthorized{Error: err.Error()})
	case http.StatusForbidden:
		return e.JSON(status, api.ErrorForbidden{Error: err.Error()})
	default:
		return e.JSON(status, map[string]string{"error": err.Error()})
	}
}
