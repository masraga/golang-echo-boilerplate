package auth

import "context"

func (r *AuthRepository) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {
	stmt := r.Sql.InsertInto(TableAuth).
		Set("phone_no", "08234567890").
		Set("pin", "123456").
		Set("otp_code", "1234")

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	return
}
