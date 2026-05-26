package auth

import "context"

func (s *AuthService) VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error) {
	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{PhoneNo: input.PhoneNo})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	verify, err := s.AuthRepositoryWriter.VerifyUserAccount(ctx, VerifyUserAccountInput{
		PhoneNo: authUser.PhoneNo,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	output.UserId = authUser.Id
	output.PhoneNo = authUser.PhoneNo
	output.IsVerified = verify.IsVerified
	output.IsNewUser = authUser.PinCode == nil
	return
}
