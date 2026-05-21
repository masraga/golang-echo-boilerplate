package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error) {
	stmt := r.Sql.InsertInto(AuthCodeTableName).
		Set("otp_code", input.OtpCode).
		Set("user_id", input.UserId).
		Set("note", input.Note).
		Set("expired_at_utc0", input.ExpiredAtUtc0)
	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = errors.Join(err, ErrCreateNewOTP)
	}
	return
}
