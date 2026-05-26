package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error) {
	stmt := r.Sql.Update(AuthCodeTableName).
		Set("is_verified", true).
		Where("user_id = ?", input.UserId).
		Where("otp_code = ?", input.OtpCode)
	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrVerifyOtp))
		return
	}
	output.IsValid = true
	output.UserId = input.UserId
	return
}
