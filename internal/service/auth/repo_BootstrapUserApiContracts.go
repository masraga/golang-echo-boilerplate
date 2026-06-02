package auth

import (
	"context"
	"errors"
)

func (r *AuthRepository) BootstrapUserApiContracts(ctx context.Context, input BootstrapUserApiContractsInput) (output BootstrapUserApiContractsOutput, err error) {
	res, err := r.UseTx(ctx).ExecContext(ctx, `
INSERT INTO public.auth_user_api_contract (
    user_id,
    api_contract_id,
    created_at_utc0,
    updated_at_utc0,
    is_active
)
SELECT
    $1,
    aac.id,
    floor(extract(epoch FROM now()) * 1000)::bigint,
    floor(extract(epoch FROM now()) * 1000)::bigint,
    TRUE
FROM public.auth_api_contract aac
WHERE aac.is_active = TRUE
AND NOT EXISTS (
    SELECT 1
    FROM public.auth_user_api_contract auac
    WHERE auac.user_id = $1
    AND auac.api_contract_id = aac.id
    AND auac.is_active = TRUE
)`, input.UserId)
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrBootstrapUserApiContract))
		return
	}
	output.InsertedCount, err = res.RowsAffected()
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrBootstrapUserApiContract))
		return
	}
	return
}
