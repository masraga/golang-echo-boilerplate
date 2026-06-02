package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) ValidateJwtToken(ctx context.Context, input ValidateJwtTokenInput) (output ValidateJwtTokenOutput, err error) {
	claims := &ValidateJwtTokenOutput{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrAuthSigInvalid
		}
		return []byte(s.JwtSecret), nil
	})

	if err != nil {
		err = s.Err.Wrap(ErrAuthSigInvalid)
		return
	}
	if !token.Valid {
		err = s.Err.Wrap(ErrAuthTokenInvalid)
		return
	}
	output = *claims

	if claims.ExpiredAtUtc0 < time.Now().UnixMilli() {
		err = s.Err.Wrap(ErrAuthTokenExpired)
		return
	}

	_, err = s.AuthRepositoryReader.FindAccessToken(ctx, FindAccessTokenInput{
		Token:  input.Token,
		UserId: claims.UserId,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	return
}
