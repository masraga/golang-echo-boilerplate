package auth

import (
	"context"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
)

func (s *AuthService) CreateAuthRoleContractApi(ctx context.Context, input CreateAuthRoleContractApiInput) (output CreateAuthRoleContractApiOutput, err error) {
	now := utiltime.NowUtc0()
	output, err = s.AuthRepositoryWriter.CreateAuthRoleContractApi(ctx, CreateAuthRoleContractApiOutput{
		RoleId:            input.RoleId,
		AuthApiContractId: input.AuthApiContractId,
		CreatedAtUtc0:     now,
		UpdatedAtUtc0:     now,
		CreatedBy:         input.CreatedBy,
		IsActive:          true,
	})
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) ListAuthRoleContractApis(ctx context.Context, input ListAuthRoleContractApisInput) (output ListAuthRoleContractApisOutput, err error) {
	output, err = s.AuthRepositoryReader.ListAuthRoleContractApis(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) DeleteAuthRoleContractApi(ctx context.Context, input DeleteAuthRoleContractApiInput) (output DeleteAuthRoleContractApiOutput, err error) {
	output, err = s.AuthRepositoryWriter.DeleteAuthRoleContractApi(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
