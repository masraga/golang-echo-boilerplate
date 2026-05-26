package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func (s *CryptoService) Decrypt(ctx context.Context, input DecryptInput) (output DecryptOutput, err error) {
	data, err := base64.StdEncoding.DecodeString(input.HashCode)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrDecodeHashString))
		return
	}

	block, err := aes.NewCipher([]byte(s.ConfigCryptoKey))
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrCreateChipher))
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrCreateChiperGCM))
		return
	}

	nonceSize := gcm.NonceSize()
	dataLen := len(data)
	if dataLen < nonceSize {
		err = s.Err.Wrap(ErrDecryptHashString)
		return
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrDecryptHashString))
		return
	}

	output.Result = string(plaintext)

	return
}
