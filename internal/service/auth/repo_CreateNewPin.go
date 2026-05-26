package auth

import (
	"context"
	"errors"
	"time"
)

func (r *AuthRepository) CreateNewPin(ctx context.Context, input CreateNewPinInput) (output CreateNewPinOutput, err error) {
	stmt := r.Sql.Update(TableAuth).
		Set("pin", input.PinCode).
		Set("updated_at_utc0", time.Now().UnixMilli()).
		Where("id = ?", input.UserId)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateNewPin))
		return
	}
	output.PinCode = input.PinCode
	output.UserId = input.UserId
	return
}
