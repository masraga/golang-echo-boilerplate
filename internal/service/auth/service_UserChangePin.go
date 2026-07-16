package auth

import (
	"context"

	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
)

func (s *AuthService) UserChangePin(ctx context.Context, input UserChangePinInput) (output UserChangePinOutput, err error) {
	user, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		UserId: input.UserId,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	if user.PinCode == nil {
		err = s.Err.Wrap(ErrPinNotDefined)
		return
	}
	if pointer.SafeString(user.PinCode) != input.OldPin {
		err = s.Err.Wrap(ErrPinCodeNotMatch)
		return
	}
	if input.NewPin != input.RetypeNewPin {
		err = s.Err.Wrap(ErrPinCodeNotMatch)
		return
	}

	updateAuth, err := s.AuthRepositoryWriter.CreateNewPin(ctx, CreateNewPinInput{
		PinCode: input.NewPin,
		UserId:  user.Id,
	})
	if err != nil {
		err = s.Err.Wrap(ErrCreateNewPin)
		return
	}

	output.IsUpdate = updateAuth.PinCode != ""

	return
}
