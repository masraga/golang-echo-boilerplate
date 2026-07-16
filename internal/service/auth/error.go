package auth

import "errors"

var (
	ErrClaimJwtToken           error = errors.New("error to claims new token")
	ErrBeginDbTx               error = errors.New("error to begin database transaction")
	ErrCreateNewAccount        error = errors.New("error to create new account")
	ErrUpdateFirebaseId        error = errors.New("error to update firebase id")
	ErrUpdateOtpValidity       error = errors.New("error to update otp validity")
	ErrAuthNotFound            error = errors.New("error auth data not found")
	ErrDuplicateUser           error = errors.New("error user already registered")
	ErrFindOTPNotFound         error = errors.New("error otp data not found")
	ErrOTPAlreadyExist         error = errors.New("error otp already exist")
	ErrFailedDeleteUserOTP     error = errors.New("error to delete user otp")
	ErrCreateNewOTP            error = errors.New("error to create new otp")
	ErrOtpIsExpired            error = errors.New("error OTP is expired")
	ErrVerifyOtp               error = errors.New("error to verify OTP")
	ErrOtpAleadyVerified       error = errors.New("error OTP already verified")
	ErrVerifyUserAccount       error = errors.New("error to verify user account")
	ErrValidateRetypePin       error = errors.New("error to validate retype pin code")
	ErrPinCodeNotMatch         error = errors.New("error pin code not match")
	ErrPinIsTooLongOrShort     error = errors.New("error pin code is too long or short")
	ErrOtpValidationRequired   error = errors.New("must validate otp before validating pin")
	ErrCreateNewPin            error = errors.New("error to create new pin code")
	ErrStoreAccessToken        error = errors.New("error to store access token")
	ErrFindAccessTokenNotFound error = errors.New("error access token data not found")
	ErrAuthSigInvalid          error = errors.New("error invalid signature")
	ErrAuthTokenInvalid        error = errors.New("error invalid token")
	ErrAuthTokenExpired        error = errors.New("error token is expired")
	ErrPinNotDefined           error = errors.New("pin not setup yet")

	ErrCreateAuthApiContract       error = errors.New("error to create auth api contract")
	ErrFindAuthApiContractNotFound error = errors.New("error auth api contract data not found")
	ErrUpdateAuthApiContract       error = errors.New("error to update auth api contract")
	ErrDeleteAuthApiContract       error = errors.New("error to delete auth api contract")

	ErrCreateAuthUserApiContract       error = errors.New("error to create auth user api contract")
	ErrFindAuthUserApiContractNotFound error = errors.New("error auth user api contract data not found")
	ErrUpdateAuthUserApiContract       error = errors.New("error to update auth user api contract")
	ErrDeleteAuthUserApiContract       error = errors.New("error to delete auth user api contract")

	ErrUserApiContractForbidden error = errors.New("error user does not have access to api contract")
	ErrBootstrapUserApiContract error = errors.New("error to bootstrap user api contract")

	ErrCreateAuthRole       error = errors.New("error to create auth role")
	ErrFindAuthRoleNotFound error = errors.New("error auth role data not found")
	ErrUpdateAuthRole       error = errors.New("error to update auth role")
	ErrDeleteAuthRole       error = errors.New("error to delete auth role")

	ErrCreateAuthRoleContractApi       error = errors.New("error to create auth role contract api")
	ErrFindAuthRoleContractApiNotFound error = errors.New("error auth role contract api data not found")
	ErrDeleteAuthRoleContractApi       error = errors.New("error to delete auth role contract api")

	ErrAssignAuthUserRole error = errors.New("error to assign auth user role")
	ErrDeleteAuthUserRole error = errors.New("error to delete auth user role")
)
