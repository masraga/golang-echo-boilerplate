package generator

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandom(length int, numbersOnly bool) (string, error) {
	charset := alphaNumericCharset

	if numbersOnly {
		charset = numberCharset
	}

	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
