package auth

import "errors"

var (
	ErrClaimJwtToken    error = errors.New("error to claims new token")
	ErrBeginDbTx        error = errors.New("error to begin database transaction")
	ErrCreateNewAccount error = errors.New("error to create new account")
	ErrAuthNotFound     error = errors.New("error auth data not found")
	ErrDuplicateUser    error = errors.New("error user already registered")
)
