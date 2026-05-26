package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) FindOTP(ctx context.Context, input FindOTPInput) (output FindOTPOutput, err error) {
	stmt := r.Sql.From(AuthCodeTableName+" ac").
		Select("ac.id").To(&output.Id).
		Select("ac.otp_code").To(&output.OtpCode).
		Select("ac.note").To(&output.Note).
		Select("ac.expired_at_utc0").To(&output.ExpiredAtUtc0).
		Select("ac.is_verified").To(&output.IsVerified).
		Where("ac.user_id = ?", input.UserId).
		Where("ac.otp_code = ?", input.OtpCode).
		Where("ac.is_active = ?", true)
	err = stmt.QueryRowAndClose(ctx, r.Db)
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindOTPNotFound))
		return
	}
	return
}
