package auth

import (
	"context"
	"database/sql"
	"errors"
)

func (r *AuthRepository) FindAuth(ctx context.Context, input FindAuthInput) (output FindAuthOutput, err error) {
	stmt := r.Sql.From(TableAuth + " ta").
		Select("ta.id").To(&output.Id).
		Select("ta.phone_no").To(&output.PhoneNo).
		Select("ta.pin").To(&output.PinCode).
		Select("ta.is_otp_valid").To(&output.IsOtpValid)

	stmt.Where("phone_no = ?", input.PhoneNo)
	if input.UserId != "" {
		stmt.Where("id = ?", input.UserId)
	}
	stmt.Where("is_active = ?", true)
	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = r.Err.Wrap(errors.Join(err, ErrAuthNotFound))
			return
		}
	}
	return
}
