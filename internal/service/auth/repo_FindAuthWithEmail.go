package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) FindAuthWithEmail(ctx context.Context, input FindAuthWithEmailInput) (output FindAuthWithEmailOutput, err error) {
	stmt := r.Sql.From(TableAuth).
		Select("id").To(&output.Id).
		Select("phone_no").To(&output.PhoneNo).
		Select("email").To(&output.Email).
		Where("is_active = ?", true).
		Where("email = ?", input.Email)
	err = stmt.QueryRowAndClose(ctx, r.Db)
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindAuthWithEmail))
		return
	}
	return
}
