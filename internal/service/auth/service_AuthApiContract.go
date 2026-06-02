package auth

import (
	"context"

	"github.com/masraga/kerp-api/internal/util/parser"
	utiltime "github.com/masraga/kerp-api/internal/util/time"
)

func (s *AuthService) CreateAuthApiContract(ctx context.Context, input CreateAuthApiContractInput) (output CreateAuthApiContractOutput, err error) {
	now := utiltime.NowUtc0()
	output, err = s.AuthRepositoryWriter.CreateAuthApiContract(ctx, CreateAuthApiContractOutput{
		Id:             input.Id,
		EndpointPath:   parser.NormalizeEndpointPath(input.EndpointPath),
		EndpointMethod: parser.NormalizeEndpointMethod(input.EndpointMethod),
		Description:    input.Description,
		CreatedAtUtc0:  now,
		UpdatedAtUtc0:  now,
		IsActive:       true,
	})
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) GetAuthApiContract(ctx context.Context, input GetAuthApiContractInput) (output GetAuthApiContractOutput, err error) {
	output, err = s.AuthRepositoryReader.GetAuthApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) ListAuthApiContracts(ctx context.Context, input ListAuthApiContractsInput) (output ListAuthApiContractsOutput, err error) {
	output, err = s.AuthRepositoryReader.ListAuthApiContracts(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) UpdateAuthApiContract(ctx context.Context, input UpdateAuthApiContractInput) (output UpdateAuthApiContractOutput, err error) {
	input.EndpointPath = parser.NormalizeEndpointPath(input.EndpointPath)
	input.EndpointMethod = parser.NormalizeEndpointMethod(input.EndpointMethod)
	output, err = s.AuthRepositoryWriter.UpdateAuthApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}

func (s *AuthService) DeleteAuthApiContract(ctx context.Context, input DeleteAuthApiContractInput) (output DeleteAuthApiContractOutput, err error) {
	output, err = s.AuthRepositoryWriter.DeleteAuthApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
