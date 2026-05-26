package auth

import "context"

func (r *AuthRepository) DeleteAllUserOTP(ctx context.Context, input DeleteAllUserOTPInput) (output DeleteAllUserOTPOutput, err error) {
	stmt := r.Sql.Update(AuthCodeTableName).
		Set(`is_active`, false).
		Where("user_id = ? AND is_active = ?", input.UserId, true)
	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	output.IsSuccess = err == nil
	if err != nil {
		err = r.Err.Wrap(err)
	}
	return
}
