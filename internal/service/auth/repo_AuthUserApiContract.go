package auth

import (
	"context"
	"database/sql"
	"errors"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
)

func (r *AuthRepository) CreateAuthUserApiContract(ctx context.Context, input CreateAuthUserApiContractOutput) (output CreateAuthUserApiContractOutput, err error) {
	stmt := r.Sql.InsertInto(AuthUserApiContractTable).
		Set("user_id", input.UserId).
		Set("api_contract_id", input.ApiContractId).
		Set("created_at_utc0", input.CreatedAtUtc0).
		Set("updated_at_utc0", input.UpdatedAtUtc0).
		Set("is_active", input.IsActive).
		Returning("id").To(&output.Id)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthUserApiContract))
		return
	}
	output.UserId = input.UserId
	output.ApiContractId = input.ApiContractId
	output.CreatedAtUtc0 = input.CreatedAtUtc0
	output.UpdatedAtUtc0 = input.UpdatedAtUtc0
	output.IsActive = input.IsActive
	return
}

func (r *AuthRepository) GetAuthUserApiContract(ctx context.Context, input GetAuthUserApiContractInput) (output GetAuthUserApiContractOutput, err error) {
	stmt := r.Sql.From(AuthUserApiContractTable+" auac").
		Select("auac.id").To(&output.Id).
		Select("auac.user_id").To(&output.UserId).
		Select("auac.api_contract_id").To(&output.ApiContractId).
		Select("auac.created_at_utc0").To(&output.CreatedAtUtc0).
		Select("auac.updated_at_utc0").To(&output.UpdatedAtUtc0).
		Select("auac.is_active").To(&output.IsActive).
		Where("auac.id = ?", input.Id).
		Where("auac.is_active = ?", true)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindAuthUserApiContractNotFound))
		return
	}
	return
}

func (r *AuthRepository) ListAuthUserApiContracts(ctx context.Context, input ListAuthUserApiContractsInput) (output ListAuthUserApiContractsOutput, err error) {
	var item AuthUserApiContract
	stmt := r.Sql.From(AuthUserApiContractTable+" auac").
		Select("auac.id").To(&item.Id).
		Select("auac.user_id").To(&item.UserId).
		Select("auac.api_contract_id").To(&item.ApiContractId).
		Select("auac.created_at_utc0").To(&item.CreatedAtUtc0).
		Select("auac.updated_at_utc0").To(&item.UpdatedAtUtc0).
		Select("auac.is_active").To(&item.IsActive).
		Where("auac.is_active = ?", true).
		OrderBy("auac.created_at_utc0 DESC")

	err = stmt.QueryAndClose(ctx, r.UseTx(ctx), func(rows *sql.Rows) {
		output.Data = append(output.Data, item)
	})
	if err != nil {
		err = r.Err.Wrap(err)
		return
	}
	return
}

func (r *AuthRepository) UpdateAuthUserApiContract(ctx context.Context, input UpdateAuthUserApiContractInput) (output UpdateAuthUserApiContractOutput, err error) {
	now := utiltime.NowUtc0()
	stmt := r.Sql.Update(AuthUserApiContractTable).
		Set("user_id", input.UserId).
		Set("api_contract_id", input.ApiContractId).
		Set("updated_at_utc0", now).
		Set("is_active", input.IsActive).
		Where("id = ?", input.Id)

	stmt.Returning("id").To(&output.Id).
		Returning("user_id").To(&output.UserId).
		Returning("api_contract_id").To(&output.ApiContractId).
		Returning("created_at_utc0").To(&output.CreatedAtUtc0).
		Returning("updated_at_utc0").To(&output.UpdatedAtUtc0).
		Returning("is_active").To(&output.IsActive)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUpdateAuthUserApiContract, ErrFindAuthUserApiContractNotFound))
		return
	}
	return
}

func (r *AuthRepository) DeleteAuthUserApiContract(ctx context.Context, input DeleteAuthUserApiContractInput) (output DeleteAuthUserApiContractOutput, err error) {
	stmt := r.Sql.DeleteFrom(AuthUserApiContractTable).
		Where("id = ?", input.Id)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthUserApiContract))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrFindAuthUserApiContractNotFound, err))
		return
	}
	output.IsSuccess = true
	return
}

func (r *AuthRepository) DeleteAuthUserApiContractsByUserId(ctx context.Context, input DeleteAuthUserApiContractsByUserIdInput) (output DeleteAuthUserApiContractsByUserIdOutput, err error) {
	stmt := r.Sql.DeleteFrom(AuthUserApiContractTable).
		Where("user_id = ?", input.UserId)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthUserApiContract))
		return
	}
	output.DeletedCount, err = res.RowsAffected()
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthUserApiContract))
		return
	}
	return
}

func (r *AuthRepository) InsertAuthUserApiContractsFromRole(ctx context.Context, input InsertAuthUserApiContractsFromRoleInput) (output InsertAuthUserApiContractsFromRoleOutput, err error) {
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
    arca.auth_api_contract_id,
    floor(extract(epoch FROM now()) * 1000)::bigint,
    floor(extract(epoch FROM now()) * 1000)::bigint,
    TRUE
FROM public.auth_roles_contract_api arca
WHERE arca.role_id = $2
AND arca.is_active = TRUE`, input.UserId, input.RoleId)
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthUserApiContract))
		return
	}
	output.InsertedCount, err = res.RowsAffected()
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthUserApiContract))
		return
	}
	return
}
