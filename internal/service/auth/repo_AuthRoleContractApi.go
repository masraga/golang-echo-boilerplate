package auth

import (
	"context"
	"database/sql"
	"errors"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
)

func (r *AuthRepository) CreateAuthRoleContractApi(ctx context.Context, input CreateAuthRoleContractApiOutput) (output CreateAuthRoleContractApiOutput, err error) {
	stmt := r.Sql.InsertInto(AuthRolesContractApiTable).
		Set("role_id", input.RoleId).
		Set("auth_api_contract_id", input.AuthApiContractId).
		Set("created_at_utc0", input.CreatedAtUtc0).
		Set("updated_at_utc0", input.UpdatedAtUtc0).
		Set("created_by", input.CreatedBy).
		Set("is_active", input.IsActive).
		Returning("id").To(&output.Id)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthRoleContractApi))
		return
	}
	output.RoleId = input.RoleId
	output.AuthApiContractId = input.AuthApiContractId
	output.CreatedAtUtc0 = input.CreatedAtUtc0
	output.UpdatedAtUtc0 = input.UpdatedAtUtc0
	output.CreatedBy = input.CreatedBy
	output.IsActive = input.IsActive
	return
}

func (r *AuthRepository) ListAuthRoleContractApis(ctx context.Context, input ListAuthRoleContractApisInput) (output ListAuthRoleContractApisOutput, err error) {
	var item AuthRoleContractApi
	stmt := r.Sql.From(AuthRolesContractApiTable+" arca").
		Select("arca.id").To(&item.Id).
		Select("arca.role_id").To(&item.RoleId).
		Select("arca.auth_api_contract_id").To(&item.AuthApiContractId).
		Select("arca.created_at_utc0").To(&item.CreatedAtUtc0).
		Select("arca.updated_at_utc0").To(&item.UpdatedAtUtc0).
		Select("arca.created_by").To(&item.CreatedBy).
		Select("arca.is_active").To(&item.IsActive).
		Where("arca.role_id = ?", input.RoleId).
		Where("arca.is_active = ?", true).
		OrderBy("arca.created_at_utc0 DESC")

	err = stmt.QueryAndClose(ctx, r.UseTx(ctx), func(rows *sql.Rows) {
		output.Data = append(output.Data, item)
	})
	if err != nil {
		err = r.Err.Wrap(err)
		return
	}
	return
}

func (r *AuthRepository) DeleteAuthRoleContractApi(ctx context.Context, input DeleteAuthRoleContractApiInput) (output DeleteAuthRoleContractApiOutput, err error) {
	stmt := r.Sql.Update(AuthRolesContractApiTable).
		Set("updated_at_utc0", utiltime.NowUtc0()).
		Set("is_active", false).
		Where("id = ?", input.Id).
		Where("role_id = ?", input.RoleId).
		Where("is_active = ?", true)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthRoleContractApi))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrFindAuthRoleContractApiNotFound, err))
		return
	}
	output.IsSuccess = true
	return
}
