package crypto

import "errors"

var (
	ErrCreateChipher     error = errors.New("error when create encryption chiper")
	ErrCreateChiperGCM   error = errors.New("error when create chiper GCM")
	ErrCreateNonce       error = errors.New("error when create encryption nonce")
	ErrDecodeHashString  error = errors.New("error when decode hash string")
	ErrDecryptHashString error = errors.New("error when decrypt hash string")
)
