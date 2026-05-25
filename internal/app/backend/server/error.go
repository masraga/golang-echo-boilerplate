package server

import (
	"net/http"

	"github.com/masraga/kerp-api/internal/service/auth"
)

var mapError = map[string]int{
	// 400 Error
	auth.ErrDuplicateUser.Error(): http.StatusBadRequest,
}
