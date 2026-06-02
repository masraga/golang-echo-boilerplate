package server

import (
	"net/http"

	"github.com/masraga/kerp-api/internal/service/auth"
)

var mapError = map[error]int{
	// 400 Error
	auth.ErrDuplicateUser: http.StatusBadRequest,
	auth.ErrAuthNotFound:  http.StatusBadRequest,

	auth.ErrCreateAuthApiContract:           http.StatusBadRequest,
	auth.ErrFindAuthApiContractNotFound:     http.StatusBadRequest,
	auth.ErrUpdateAuthApiContract:           http.StatusBadRequest,
	auth.ErrDeleteAuthApiContract:           http.StatusBadRequest,
	auth.ErrCreateAuthUserApiContract:       http.StatusBadRequest,
	auth.ErrFindAuthUserApiContractNotFound: http.StatusBadRequest,
	auth.ErrUpdateAuthUserApiContract:       http.StatusBadRequest,
	auth.ErrDeleteAuthUserApiContract:       http.StatusBadRequest,
	auth.ErrBootstrapUserApiContract:        http.StatusBadRequest,
	auth.ErrUserApiContractForbidden:        http.StatusForbidden,
	auth.ErrCreateAuthRole:                  http.StatusBadRequest,
	auth.ErrFindAuthRoleNotFound:            http.StatusBadRequest,
	auth.ErrUpdateAuthRole:                  http.StatusBadRequest,
	auth.ErrDeleteAuthRole:                  http.StatusBadRequest,
	auth.ErrCreateAuthRoleContractApi:       http.StatusBadRequest,
	auth.ErrFindAuthRoleContractApiNotFound: http.StatusBadRequest,
	auth.ErrDeleteAuthRoleContractApi:       http.StatusBadRequest,
	auth.ErrAssignAuthUserRole:              http.StatusBadRequest,
	auth.ErrDeleteAuthUserRole:              http.StatusBadRequest,
}
