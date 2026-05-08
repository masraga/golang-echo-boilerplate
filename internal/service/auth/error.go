package auth

import "errors"

var (
	ErrClaimJwtToken error = errors.New("error to claims new token")
	ErrBeginDbTx     error = errors.New("error to begin database transaction")
)
