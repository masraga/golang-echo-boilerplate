package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {
	stmt := r.Sql.InsertInto(TableAuth).
		Set("id", input.Id).
		Set("phone_no", input.PhoneNo)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = errors.Join(err, ErrCreateNewAccount)
		return
	}
	output.Id = input.PhoneNo
	if err != nil {
		err = r.Err.Wrap(err)
		return
	}
	return
}
