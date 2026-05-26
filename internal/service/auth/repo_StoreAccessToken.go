package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) StoreAccessToken(ctx context.Context, input StoreAccessTokenInput) (output StoreAccessTokenOutput, err error) {
	deactivateStmt := r.Sql.Update(AccessTokenTable).
		Set("is_active", false).
		Where("user_id = ? AND is_active = ?", input.UserId, true)

	_, err = deactivateStmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrStoreAccessToken))
		return
	}

	insertStmt := r.Sql.InsertInto(AccessTokenTable).
		Set("id", input.Token).
		Set("user_id", input.UserId).
		Set("expired_at_utc0", input.ExpiredAtUtc0).
		Set("is_active", true)

	_, err = insertStmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrStoreAccessToken))
		return
	}

	output.Token = input.Token
	output.UserId = input.UserId
	output.ExpiredAtUtc0 = input.ExpiredAtUtc0
	output.IsActive = true
	return
}
