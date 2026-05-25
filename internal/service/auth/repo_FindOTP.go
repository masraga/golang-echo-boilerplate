package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) FindOTP(ctx context.Context, input FindOTPInput) (output FindOTPOutput, err error) {
	stmt := r.Sql.From(AuthCodeTableName+" ac").
		Select("ac.id").To(&output.Id).
		Select("ac.otp_code").To(&output.OtpCode).
		Where("ac.user_id = ?", input.UserId).
		Where("ac.is_active = ?", true)
	err = stmt.QueryRowAndClose(ctx, r.Db)
	if err != nil {
		err = errors.Join(err, ErrFindOTPNotFound)
		return
	}
	return
}
