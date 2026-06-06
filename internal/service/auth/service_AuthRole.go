package auth

import (
	"context"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
)

func (s *AuthService) CreateAuthRole(ctx context.Context, input CreateAuthRoleInput) (output CreateAuthRoleOutput, err error) {
	now := utiltime.NowUtc0()
	output, err = s.AuthRepositoryWriter.CreateAuthRole(ctx, CreateAuthRoleOutput{
		RoleName:      input.RoleName,
		Description:   input.Description,
		OwnerId:       input.OwnerId,
		CreatedAtUtc0: now,
		UpdatedAtUtc0: now,
		CreatedBy:     input.CreatedBy,
		IsActive:      true,
	})
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) GetAuthRole(ctx context.Context, input GetAuthRoleInput) (output GetAuthRoleOutput, err error) {
	output, err = s.AuthRepositoryReader.GetAuthRole(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) ListAuthRoles(ctx context.Context, input ListAuthRolesInput) (output ListAuthRolesOutput, err error) {
	output, err = s.AuthRepositoryReader.ListAuthRoles(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) UpdateAuthRole(ctx context.Context, input UpdateAuthRoleInput) (output UpdateAuthRoleOutput, err error) {
	output, err = s.AuthRepositoryWriter.UpdateAuthRole(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) DeleteAuthRole(ctx context.Context, input DeleteAuthRoleInput) (output DeleteAuthRoleOutput, err error) {
	output, err = s.AuthRepositoryWriter.DeleteAuthRole(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
