package auth

import (
	"context"
	"errors"
	"time"
)

func (r *AuthRepository) UpdateFirebaseId(ctx context.Context, input UpdateFirebaseIdInput) (output UpdateFirebaseIdOutput, err error) {
	stmt := r.Sql.Update(TableAuth).
		Set("firebase_id", input.FirebaseId).
		Set("updated_at_utc0", time.Now().UnixMilli()).
		Where("id = ?", input.UserId)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUpdateFirebaseId))
		return
	}
	output.UserId = input.UserId
	output.FirebaseId = input.FirebaseId
	return
}
