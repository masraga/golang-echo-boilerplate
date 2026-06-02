package auth

import "context"

func (s *AuthService) BootstrapUserApiContracts(ctx context.Context, input BootstrapUserApiContractsInput) (output BootstrapUserApiContractsOutput, err error) {
	output, err = s.AuthRepositoryWriter.BootstrapUserApiContracts(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
	}
	return
}
