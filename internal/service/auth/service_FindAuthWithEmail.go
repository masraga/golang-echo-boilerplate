package auth

import "context"

func (s *AuthService) FindAuthWithEmail(ctx context.Context, input FindAuthWithEmailInput) (output FindAuthWithEmailOutput, err error) {
	output, err = s.AuthRepositoryReader.FindAuthWithEmail(ctx, input)
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}
	return
}
