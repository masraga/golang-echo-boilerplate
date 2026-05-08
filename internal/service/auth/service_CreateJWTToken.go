package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/masraga/kerp-api/internal/util/traceerr"
)

func (s *AuthService) CreateJWTToken(ctx context.Context, input CreateJWTTokenInput) (output CreateJWTTokenOutput, err error) {
	defer traceerr.WrapReturn(&err)

	invalidExpired := input.ExpiredAtUtc0 == 0
	invalidUserId := input.UserId == ""

	if invalidExpired || invalidUserId {
		return output, ErrClaimJwtToken
	}

	claims := jwt.MapClaims{
		"userId": input.UserId,
		"exp":    input.ExpiredAtUtc0,
		"iat":    time.Now().UnixMilli(),
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
