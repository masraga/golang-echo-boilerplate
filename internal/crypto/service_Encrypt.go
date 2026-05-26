package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func (s *CryptoService) Encrypt(ctx context.Context, input EncryptInput) (output EncryptOutput, err error) {
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrCreateNonce))
		return
	}

	chipertext := gcm.Seal(nonce, nonce, []byte(input.PlainText), nil)
	encode := base64.StdEncoding.EncodeToString(chipertext)
	output.Result = encode
	return
}
