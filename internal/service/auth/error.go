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
	ErrOtpIsExpired        error = errors.New("error OTP is expired")
	ErrVerifyOtp           error = errors.New("error to verify OTP")
	ErrOtpAleadyVerified   error = errors.New("error OTP already verified")
	ErrVerifyUserAccount   error = errors.New("error to verify user account")
	ErrValidateRetypePin   error = errors.New("error to validate retype pin code")
	ErrPinCodeNotMatch     error = errors.New("error pin code not match")
	ErrPinIsTooLongOrShort error = errors.New("error pin code is too long or short")
	ErrCreateNewPin        error = errors.New("error to create new pin code")
	ErrStoreAccessToken    error = errors.New("error to store access token")
)
