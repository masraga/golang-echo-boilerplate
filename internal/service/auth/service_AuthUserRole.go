package auth

import (
	"context"
	"errors"
)

func (s *AuthService) AssignAuthUserRole(ctx context.Context, input AssignAuthUserRoleInput) (output AssignAuthUserRoleOutput, err error) {
	role, err := s.AuthRepositoryReader.GetAuthRole(ctx, GetAuthRoleInput{Id: input.RoleId})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}
	input.RoleName = role.RoleName

	ctx, err = s.AuthRepositoryWriter.Begin(ctx, nil)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrBeginDbTx))
		return
	}
	defer func() {
		commitErr := s.AuthRepositoryWriter.CommitOrRollback(ctx, err)
		err = s.Err.Wrap(commitErr)
	}()

	output, err = s.AuthRepositoryWriter.AssignAuthUserRole(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	_, err = s.AuthRepositoryWriter.DeleteAuthUserApiContractsByUserId(ctx, DeleteAuthUserApiContractsByUserIdInput{
		UserId: input.UserId,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	grants, err := s.AuthRepositoryWriter.InsertAuthUserApiContractsFromRole(ctx, InsertAuthUserApiContractsFromRoleInput{
		UserId:    input.UserId,
		RoleId:    input.RoleId,
		CreatedBy: input.CreatedBy,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}
	output.GrantedCount = grants.InsertedCount
	return
}

func (s *AuthService) DeleteAuthUserRole(ctx context.Context, input DeleteAuthUserRoleInput) (output DeleteAuthUserRoleOutput, err error) {
	ctx, err = s.AuthRepositoryWriter.Begin(ctx, nil)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrBeginDbTx))
		return
	}
	defer func() {
		commitErr := s.AuthRepositoryWriter.CommitOrRollback(ctx, err)
		err = s.Err.Wrap(commitErr)
	}()

	output, err = s.AuthRepositoryWriter.DeleteAuthUserRole(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	_, err = s.AuthRepositoryWriter.DeleteAuthUserApiContractsByUserId(ctx, DeleteAuthUserApiContractsByUserIdInput{
		UserId: input.UserId,
	})
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
