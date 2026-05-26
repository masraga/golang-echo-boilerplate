package crypto_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCryptoService_Decrypt(t *testing.T) {
	const (
		validKey  = crypto.ConfigCryptoKey("0123456789abcdef")
		plainText = "sensitive value"
	)

	svc := crypto.NewCryptoService(crypto.CryptoServiceOpts{
		ConfigCryptoKey: validKey,
		Err: ctxerr.NewCtxErr(ctxerr.CtxErrOpts{
			Logger: zerolog.Nop(),
		}),
	})
	encrypted, err := svc.Encrypt(context.Background(), crypto.EncryptInput{
		PlainText: plainText,
	})
	require.NoError(t, err)

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Result)
	require.NoError(t, err)
	ciphertext[len(ciphertext)-1] ^= 0xff
	tamperedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	tests := []struct {
		name        string
		key         crypto.ConfigCryptoKey
		hashCode    string
		expected    string
		expectedErr error
	}{
		{
			name:     "should decrypt encrypted value",
			key:      validKey,
			hashCode: encrypted.Result,
			expected: plainText,
		},
		{
			name:        "should fail when hash code is not base64",
			key:         validKey,
			hashCode:    "not-base64%",
			expectedErr: crypto.ErrDecodeHashString,
		},
		{
			name:        "should fail when decoded value is shorter than nonce",
			key:         validKey,
			hashCode:    base64.StdEncoding.EncodeToString([]byte("short")),
			expectedErr: crypto.ErrDecryptHashString,
		},
		{
			name:        "should fail when ciphertext is modified",
			key:         validKey,
			hashCode:    tamperedCiphertext,
			expectedErr: crypto.ErrDecryptHashString,
		},
		{
			name:        "should fail when crypto key is invalid",
			key:         crypto.ConfigCryptoKey("invalid-key"),
			hashCode:    encrypted.Result,
			expectedErr: crypto.ErrCreateChipher,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := crypto.NewCryptoService(crypto.CryptoServiceOpts{
				ConfigCryptoKey: tt.key,
				Err: ctxerr.NewCtxErr(ctxerr.CtxErrOpts{
					Logger: zerolog.Nop(),
				}),
			})

			output, err := svc.Decrypt(context.Background(), crypto.DecryptInput{
				HashCode: tt.hashCode,
			})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Empty(t, output.Result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, output.Result)
		})
	}
}
