package crypto

import "context"

type CryptoServiceInterface interface {
	Encrypt(ctx context.Context, input EncryptInput) (output EncryptOutput, err error)
	Decrypt(ctx context.Context, input DecryptInput) (output DecryptOutput, err error)
}
