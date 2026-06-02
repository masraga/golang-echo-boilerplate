package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) ValidateUserApiContract(ctx context.Context, input ValidateUserApiContractInput) (output ValidateUserApiContractOutput, err error) {
	var grantId string
	stmt := r.Sql.From(AuthUserApiContractTable+" auac").
		Join(AuthApiContractTable+" aac", "aac.id = auac.api_contract_id").
		Select("auac.id").To(&grantId).
		Where("auac.user_id = ?", input.UserId).
		Where("aac.endpoint_path = ?", input.EndpointPath).
		Where("aac.endpoint_method = ?", input.EndpointMethod).
		Where("auac.is_active = ?", true).
		Where("aac.is_active = ?", true)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUserApiContractForbidden))
		return
	}
	output.IsAllowed = true
	return
}
