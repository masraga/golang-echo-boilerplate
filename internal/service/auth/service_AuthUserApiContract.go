package auth

import (
	"context"

	utiltime "github.com/masraga/kerp-api/internal/util/time"
)

func (s *AuthService) CreateAuthUserApiContract(ctx context.Context, input CreateAuthUserApiContractInput) (output CreateAuthUserApiContractOutput, err error) {
	now := utiltime.NowUtc0()
	output, err = s.AuthRepositoryWriter.CreateAuthUserApiContract(ctx, CreateAuthUserApiContractOutput{
		UserId:        input.UserId,
		ApiContractId: input.ApiContractId,
		CreatedAtUtc0: now,
		UpdatedAtUtc0: now,
		IsActive:      true,
	})
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) GetAuthUserApiContract(ctx context.Context, input GetAuthUserApiContractInput) (output GetAuthUserApiContractOutput, err error) {
	output, err = s.AuthRepositoryReader.GetAuthUserApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) ListAuthUserApiContracts(ctx context.Context, input ListAuthUserApiContractsInput) (output ListAuthUserApiContractsOutput, err error) {
	output, err = s.AuthRepositoryReader.ListAuthUserApiContracts(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) UpdateAuthUserApiContract(ctx context.Context, input UpdateAuthUserApiContractInput) (output UpdateAuthUserApiContractOutput, err error) {
	output, err = s.AuthRepositoryWriter.UpdateAuthUserApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) DeleteAuthUserApiContract(ctx context.Context, input DeleteAuthUserApiContractInput) (output DeleteAuthUserApiContractOutput, err error) {
	output, err = s.AuthRepositoryWriter.DeleteAuthUserApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
