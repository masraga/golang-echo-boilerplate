package crypto_test

import (
	"context"
	"testing"

	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCryptoService_Encrypt(t *testing.T) {
	tests := []struct {
		name        string
		key         crypto.ConfigCryptoKey
		plainText   string
		expectedErr error
	}{
		{
			name:      "should encrypt plaintext",
			key:       crypto.ConfigCryptoKey("0123456789abcdef"),
			plainText: "sensitive value",
		},
		{
			name:        "should fail when crypto key is invalid",
			key:         crypto.ConfigCryptoKey("invalid-key"),
			plainText:   "sensitive value",
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

			output, err := svc.Encrypt(context.Background(), crypto.EncryptInput{
				PlainText: tt.plainText,
			})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Empty(t, output.Result)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, output.Result)

			decrypted, err := svc.Decrypt(context.Background(), crypto.DecryptInput{
				HashCode: output.Result,
			})
			require.NoError(t, err)
			require.Equal(t, tt.plainText, decrypted.Result)
		})
	}
}
