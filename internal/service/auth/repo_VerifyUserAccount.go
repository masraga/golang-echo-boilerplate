package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error) {
	stmt := r.Sql.Update(TableAuth).
		Set("is_verified", true).
		Where("phone_no = ?", input.PhoneNo).
		Where("is_active = ?", true)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrVerifyUserAccount))
		return
	}
	output.PhoneNo = input.PhoneNo
	output.IsVerified = true
	return
}
