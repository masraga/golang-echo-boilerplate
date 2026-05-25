package auth

import "errors"

var (
	ErrClaimJwtToken       error = errors.New("error to claims new token")
	ErrBeginDbTx           error = errors.New("error to begin database transaction")
	ErrCreateNewAccount    error = errors.New("error to create new account")
	ErrAuthNotFound        error = errors.New("error auth data not found")
	ErrDuplicateUser       error = errors.New("error user already registered")
	ErrFindOTPNotFound     error = errors.New("error otp data not found")
	ErrOTPAlreadyExist     error = errors.New("error otp already exist")
	ErrFailedDeleteUserOTP error = errors.New("error to delete user otp")
	ErrCreateNewOTP        error = errors.New("error to create new otp")
)
