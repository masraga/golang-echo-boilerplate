package oauth

import (
	"context"
	"errors"
)

func (r *OAuthRepository) CreateAccessToken(ctx context.Context, input CreateAccessTokenInput) (output CreateAccessTokenOutput, err error) {
	stmt := r.sql.InsertInto(OauthTableName).
		Set("id", input.Id).
		Set("access_token", input.AccessToken).
		Set("refresh_token", input.RefreshToken)
	_, err = stmt.ExecAndClose(ctx, r.db)
	output.Id = input.Id
	if err != nil {
		err = r.err.Wrap(errors.Join(err, ErrCreateNewAccessToken))
		return
	}
	return
}
