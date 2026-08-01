package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) CreateJWTToken(ctx context.Context, input CreateJWTTokenInput) (output CreateJWTTokenOutput, err error) {

	invalidExpired := input.ExpiredAtUtc0 == 0
	invalidUserId := input.UserId == ""

	if invalidExpired || invalidUserId {
		return output, ErrClaimJwtToken
	}

	issuerAtUtc0 := input.IssuerAtUtc0
	if issuerAtUtc0 == 0 {
		issuerAtUtc0 = time.Now().UnixMilli()
	}

	//set metadata to json, so it can be store to jwt
	var metadata []byte
	metadata, _ = json.Marshal(input.Metadata)

	claims := jwt.MapClaims{
		"userId":   input.UserId,
		"exp":      input.ExpiredAtUtc0,
		"iat":      issuerAtUtc0,
		"metadata": string(metadata),
	}
	newClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := newClaims.SignedString([]byte(s.JwtSecret))
	if err != nil {
		return output, ErrClaimJwtToken
	}
	output.UserId = input.UserId
	output.Token = token
	return
}
