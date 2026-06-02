package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) FindAccessToken(ctx context.Context, input FindAccessTokenInput) (output FindAccessTokenOutput, err error) {
	stmt := r.Sql.From(AccessTokenTable+" aat").
		Select("aat.id").To(&output.Token).
		Select("aat.user_id").To(&output.UserId).
		Select("aat.expired_at_utc0").To(&output.ExpiredAtUtc0).
		Select("aat.is_active").To(&output.IsActive).
		Where("aat.id = ?", input.Token).
		Where("aat.user_id = ?", input.UserId).
		Where("aat.is_active = ?", true)

	err = stmt.QueryRowAndClose(ctx, r.Db)
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindAccessTokenNotFound))
		return
	}
	return
}
