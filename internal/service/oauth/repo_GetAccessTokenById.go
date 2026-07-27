package oauth

import (
	"context"
	"errors"
)

func (r *OAuthRepository) GetAccessTokenById(ctx context.Context, input GetAccessTokenByIdInput) (output GetAccessTokenByIdOutput, err error) {
	stmt := r.sql.From(OauthTableName).
		Select("id").To(&output.Id).
		Select("username").To(&output.Username).
		Select("email").To(&output.Email).
		Select("access_token").To(&output.AccessToken).
		Select("refresh_token").To(&output.RefreshToken).
		Select("expiry_at_utc0").To(&output.ExpiryAtUtc0).
		Where("is_active = ?", true).
		Where("access_token = ?", input.AccessToken)
	err = stmt.QueryRowAndClose(ctx, r.db)
	if err != nil {
		err = r.err.Wrap(errors.Join(err, ErrAccessTokenNotFound))
		return
	}
	return
}
