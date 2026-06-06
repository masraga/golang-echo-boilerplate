package auth

import (
	"context"
	"database/sql"
	"errors"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
)

func (r *AuthRepository) CreateAuthRole(ctx context.Context, input CreateAuthRoleOutput) (output CreateAuthRoleOutput, err error) {
	stmt := r.Sql.InsertInto(AuthRolesTable).
		Set("role_name", input.RoleName).
		Set("description", input.Description).
		Set("owner_id", input.OwnerId).
		Set("created_at_utc0", input.CreatedAtUtc0).
		Set("updated_at_utc0", input.UpdatedAtUtc0).
		Set("created_by", input.CreatedBy).
		Set("is_active", input.IsActive).
		Returning("id").To(&output.Id)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrCreateAuthRole))
		return
	}
	output.RoleName = input.RoleName
	output.Description = input.Description
	output.OwnerId = input.OwnerId
	output.CreatedAtUtc0 = input.CreatedAtUtc0
	output.UpdatedAtUtc0 = input.UpdatedAtUtc0
	output.CreatedBy = input.CreatedBy
	output.IsActive = input.IsActive
	return
}

func (r *AuthRepository) GetAuthRole(ctx context.Context, input GetAuthRoleInput) (output GetAuthRoleOutput, err error) {
	stmt := r.Sql.From(AuthRolesTable+" ar").
		Select("ar.id").To(&output.Id).
		Select("ar.role_name").To(&output.RoleName).
		Select("ar.description").To(&output.Description).
		Select("ar.owner_id").To(&output.OwnerId).
		Select("ar.created_at_utc0").To(&output.CreatedAtUtc0).
		Select("ar.updated_at_utc0").To(&output.UpdatedAtUtc0).
		Select("ar.created_by").To(&output.CreatedBy).
		Select("ar.is_active").To(&output.IsActive).
		Where("ar.id = ?", input.Id).
		Where("ar.is_active = ?", true)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrFindAuthRoleNotFound))
		return
	}
	return
}

func (r *AuthRepository) ListAuthRoles(ctx context.Context, input ListAuthRolesInput) (output ListAuthRolesOutput, err error) {
	var item AuthRole
	stmt := r.Sql.From(AuthRolesTable+" ar").
		Select("ar.id").To(&item.Id).
		Select("ar.role_name").To(&item.RoleName).
		Select("ar.description").To(&item.Description).
		Select("ar.owner_id").To(&item.OwnerId).
		Select("ar.created_at_utc0").To(&item.CreatedAtUtc0).
		Select("ar.updated_at_utc0").To(&item.UpdatedAtUtc0).
		Select("ar.created_by").To(&item.CreatedBy).
		Select("ar.is_active").To(&item.IsActive).
		Where("ar.is_active = ?", true).
		OrderBy("ar.role_name ASC")

	err = stmt.QueryAndClose(ctx, r.UseTx(ctx), func(rows *sql.Rows) {
		output.Data = append(output.Data, item)
	})
	if err != nil {
		err = r.Err.Wrap(err)
		return
	}
	return
}

func (r *AuthRepository) UpdateAuthRole(ctx context.Context, input UpdateAuthRoleInput) (output UpdateAuthRoleOutput, err error) {
	stmt := r.Sql.Update(AuthRolesTable).
		Set("role_name", input.RoleName).
		Set("description", input.Description).
		Set("owner_id", input.OwnerId).
		Set("updated_at_utc0", utiltime.NowUtc0()).
		Set("created_by", input.CreatedBy).
		Set("is_active", input.IsActive).
		Where("id = ?", input.Id)

	stmt.Returning("id").To(&output.Id).
		Returning("role_name").To(&output.RoleName).
		Returning("description").To(&output.Description).
		Returning("owner_id").To(&output.OwnerId).
		Returning("created_at_utc0").To(&output.CreatedAtUtc0).
		Returning("updated_at_utc0").To(&output.UpdatedAtUtc0).
		Returning("created_by").To(&output.CreatedBy).
		Returning("is_active").To(&output.IsActive)

	err = stmt.QueryRowAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrUpdateAuthRole, ErrFindAuthRoleNotFound))
		return
	}
	return
}

func (r *AuthRepository) DeleteAuthRole(ctx context.Context, input DeleteAuthRoleInput) (output DeleteAuthRoleOutput, err error) {
	stmt := r.Sql.Update(AuthRolesTable).
		Set("updated_at_utc0", utiltime.NowUtc0()).
		Set("is_active", false).
		Where("id = ?", input.Id).
		Where("is_active = ?", true)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthRole))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrFindAuthRoleNotFound, err))
		return
	}
	output.IsSuccess = true
	return
}
