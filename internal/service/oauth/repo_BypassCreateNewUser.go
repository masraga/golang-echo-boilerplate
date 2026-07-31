package oauth

import (
	"context"
	"errors"

	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (r *OAuthRepository) BypassCreateNewUser(ctx context.Context, input BypassCreateNewUserInput) (output BypassCreateNewUserOutput, err error) {
	stmt := r.sql.InsertInto(auth.TableAuth).
		Set("id", input.Id).
		Set("phone_no", input.PhoneNo).
		Set("pin", input.Pin).
		Set("is_verified", true).
		Set("email", input.Email).
		Set("created_by", input.CreatedBy).
		Set("login_mode", input.LoginMode)
	_, err = stmt.ExecAndClose(ctx, r.db)
	if err != nil {
		err = r.err.Wrap(errors.Join(err, ErrBypassCreateUser))
		return
	}
	output.Id = input.Id
	return
}
