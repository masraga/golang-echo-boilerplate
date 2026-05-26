package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func (s *CryptoService) Decrypt(ctx context.Context, input DEncryptInput) (output DecryptOutput, err error) {
	data, err := base64.StdEncoding.DecodeString(input.HashCode)
	if err != nil {
		err = s.Err.Wrap(ErrDecodeHashString)
		return
	}

	block, err := aes.NewCipher([]byte(s.ConfigCryptoKey))
	if err != nil {
		err = s.Err.Wrap(ErrCreateChipher)
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err = s.Err.Wrap(ErrCreateChiperGCM)
		return
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		err = s.Err.Wrap(ErrDecodeHashString)
		return
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		err = s.Err.Wrap(ErrDecodeHashString)
		return
	}

	output.Result = string(plaintext)

	return
}
