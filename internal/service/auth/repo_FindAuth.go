package auth

import (
	"context"
	"database/sql"
	"errors"
)

func (r *AuthRepository) FindAuth(ctx context.Context, input FindAuthInput) (output FindAuthOutput, err error) {
	stmt := r.Sql.From(TableAuth + " ta").
		Select("ta.id").To(&output.Id).
		Select("ta.phone_no").To(&output.PhoneNo)

	if input.PhoneNo != "" {
		stmt.Where("phone_no", input.PhoneNo)
	}
	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrAuthNotFound
		}
		return
	}
	return
}
