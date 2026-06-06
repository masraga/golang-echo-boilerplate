package auth

import (
	"context"

	"github.com/masraga/golang-echo-boilerplate/internal/util/parser"
)

func (s *AuthService) ValidateUserApiContract(ctx context.Context, input ValidateUserApiContractInput) (output ValidateUserApiContractOutput, err error) {
	if s.AuthAccessBootstrapUserId != "" && input.UserId == string(s.AuthAccessBootstrapUserId) {
		output.IsAllowed = true
		return
	}

	input.EndpointPath = parser.NormalizeEndpointPath(input.EndpointPath)
	input.EndpointMethod = parser.NormalizeEndpointMethod(input.EndpointMethod)
	output, err = s.AuthRepositoryReader.ValidateUserApiContract(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
