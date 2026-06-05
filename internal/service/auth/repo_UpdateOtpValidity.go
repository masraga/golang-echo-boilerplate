package auth

import (
	"context"
	"errors"
	"time"
)

func (r *AuthRepository) UpdateOtpValidity(ctx context.Context, input UpdateOtpValidityInput) (output UpdateOtpValidityOutput, err error) {
	stmt := r.Sql.Update(TableAuth).
		Set("is_otp_valid", input.IsOtpValid).
		Set("updated_at_utc0", time.Now().UnixMilli()).
		Where("id = ?", input.UserId)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUpdateOtpValidity))
		return
	}
	output.UserId = input.UserId
	output.IsOtpValid = input.IsOtpValid
	return
}
