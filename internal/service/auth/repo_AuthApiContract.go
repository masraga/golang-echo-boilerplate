package auth

import (
	"context"
	"database/sql"
	"errors"

	utiltime "github.com/masraga/kerp-api/internal/util/time"
)

func (r *AuthRepository) CreateAuthApiContract(ctx context.Context, input CreateAuthApiContractOutput) (output CreateAuthApiContractOutput, err error) {
	stmt := r.Sql.InsertInto(AuthApiContractTable).
		Set("id", input.Id).
		Set("endpoint_path", input.EndpointPath).
		Set("endpoint_method", input.EndpointMethod).
		Set("description", input.Description).
		Set("created_at_utc0", input.CreatedAtUtc0).
		Set("updated_at_utc0", input.UpdatedAtUtc0).
		Set("is_active", input.IsActive)

	_, err = stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthApiContract))
		return
	}
	output = input
	return
}

func (r *AuthRepository) GetAuthApiContract(ctx context.Context, input GetAuthApiContractInput) (output GetAuthApiContractOutput, err error) {
	stmt := r.Sql.From(AuthApiContractTable+" aac").
		Select("aac.id").To(&output.Id).
		Select("aac.endpoint_path").To(&output.EndpointPath).
		Select("aac.endpoint_method").To(&output.EndpointMethod).
		Select("aac.description").To(&output.Description).
		Select("aac.created_at_utc0").To(&output.CreatedAtUtc0).
		Select("aac.updated_at_utc0").To(&output.UpdatedAtUtc0).
		Select("aac.is_active").To(&output.IsActive).
		Where("aac.id = ?", input.Id).
		Where("aac.is_active = ?", true)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindAuthApiContractNotFound))
		return
	}
	return
}

func (r *AuthRepository) ListAuthApiContracts(ctx context.Context, input ListAuthApiContractsInput) (output ListAuthApiContractsOutput, err error) {
	var item AuthApiContract
	stmt := r.Sql.From(AuthApiContractTable+" aac").
		Select("aac.id").To(&item.Id).
		Select("aac.endpoint_path").To(&item.EndpointPath).
		Select("aac.endpoint_method").To(&item.EndpointMethod).
		Select("aac.description").To(&item.Description).
		Select("aac.created_at_utc0").To(&item.CreatedAtUtc0).
		Select("aac.updated_at_utc0").To(&item.UpdatedAtUtc0).
		Select("aac.is_active").To(&item.IsActive).
		Where("aac.is_active = ?", true).
		OrderBy("aac.id ASC")

	err = stmt.QueryAndClose(ctx, r.UseTx(ctx), func(rows *sql.Rows) {
		output.Data = append(output.Data, item)
	})
	if err != nil {
		err = r.Err.Wrap(err)
		return
	}
	return
}

func (r *AuthRepository) UpdateAuthApiContract(ctx context.Context, input UpdateAuthApiContractInput) (output UpdateAuthApiContractOutput, err error) {
	now := utiltime.NowUtc0()
	stmt := r.Sql.Update(AuthApiContractTable).
		Set("endpoint_path", input.EndpointPath).
		Set("endpoint_method", input.EndpointMethod).
		Set("description", input.Description).
		Set("updated_at_utc0", now).
		Set("is_active", input.IsActive).
		Where("id = ?", input.Id)

	stmt.Returning("id").To(&output.Id).
		Returning("endpoint_path").To(&output.EndpointPath).
		Returning("endpoint_method").To(&output.EndpointMethod).
		Returning("description").To(&output.Description).
		Returning("created_at_utc0").To(&output.CreatedAtUtc0).
		Returning("updated_at_utc0").To(&output.UpdatedAtUtc0).
		Returning("is_active").To(&output.IsActive)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUpdateAuthApiContract, ErrFindAuthApiContractNotFound))
		return
	}
	return
}

func (r *AuthRepository) DeleteAuthApiContract(ctx context.Context, input DeleteAuthApiContractInput) (output DeleteAuthApiContractOutput, err error) {
	stmt := r.Sql.Update(AuthApiContractTable).
		Set("updated_at_utc0", utiltime.NowUtc0()).
		Set("is_active", false).
		Where("id = ?", input.Id).
		Where("is_active = ?", true)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthApiContract))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrFindAuthApiContractNotFound, err))
		return
	}
	output.IsSuccess = true
	return
}
