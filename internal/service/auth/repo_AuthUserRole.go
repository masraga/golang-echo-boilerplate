package auth

import (
	"context"
	"errors"

	utiltime "github.com/masraga/kerp-api/internal/util/time"
)

func (r *AuthRepository) AssignAuthUserRole(ctx context.Context, input AssignAuthUserRoleInput) (output AssignAuthUserRoleOutput, err error) {
	now := utiltime.NowUtc0()
	stmt := r.Sql.Update(TableAuth).
		Set("role_id", input.RoleId).
		Set("role_name", input.RoleName).
		Set("created_by", input.CreatedBy).
		Set("updated_at_utc0", now).
		Where("id = ?", input.UserId).
		Where("is_active = ?", true)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrAssignAuthUserRole))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrAuthNotFound, err))
		return
	}
	output.UserId = input.UserId
	output.RoleId = input.RoleId
	output.RoleName = input.RoleName
	output.UpdatedAtUtc0 = now
	output.CreatedBy = input.CreatedBy
	return
}

func (r *AuthRepository) DeleteAuthUserRole(ctx context.Context, input DeleteAuthUserRoleInput) (output DeleteAuthUserRoleOutput, err error) {
	stmt := r.Sql.Update(TableAuth).
		Set("role_id", nil).
		Set("role_name", nil).
		Set("created_by", input.CreatedBy).
		Set("updated_at_utc0", utiltime.NowUtc0()).
		Where("id = ?", input.UserId).
		Where("is_active = ?", true)

	res, err := stmt.ExecAndClose(ctx, r.UseTx(ctx))
	if err != nil {
		err = r.Err.Wrap(errors.Join(err, ErrDeleteAuthUserRole))
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		err = r.Err.Wrap(errors.Join(ErrAuthNotFound, err))
		return
	}
	output.IsSuccess = true
	return
}
