package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestAuthService_UserChangePin(t *testing.T) {
	var (
		validUserId       string = faker.UUIDHyphenated()
		validUserPhoneNo  string = "081234567890"
		validOldPin       string = "654321"
		validNewPin       string = "123456"
		validRetypeNewPin string = "123456"
	)

	type args struct {
		ctx   context.Context
		input auth.UserChangePinInput
	}

	type fields struct {
		AuthRepositoryReader auth.AuthRepositoryReaderInterface
		AuthRepositoryWriter auth.AuthRepositoryWriterInterface
	}

	type expected = testutil.Result[auth.UserChangePinOutput]

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success change pin",
			args: args{
				ctx: context.Background(),
				input: auth.UserChangePinInput{
					UserPhoneNo:  validUserPhoneNo,
					OldPin:       validOldPin,
					NewPin:       validNewPin,
					RetypeNewPin: validRetypeNewPin,
					AuthUserId:   validUserId,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.UserChangePinOutput{
					IsUpdate: true,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{Id: validUserPhoneNo, PinCode: &validOldPin}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					CreateNewPin(gomock.Any(), gomock.Any()).
					Return(auth.CreateNewPinOutput{PinCode: validNewPin, UserId: validUserPhoneNo}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "failed user pin not setup yet",
			args: args{
				ctx: context.Background(),
				input: auth.UserChangePinInput{
					UserPhoneNo:  validUserPhoneNo,
					OldPin:       validOldPin,
					NewPin:       validNewPin,
					RetypeNewPin: validRetypeNewPin,
					AuthUserId:   validUserId,
				},
			},
			expected: expected{
				Err:   auth.ErrPinNotDefined,
				Value: auth.UserChangePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{Id: validUserPhoneNo, PinCode: nil}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "failed user old pin invalid",
			args: args{
				ctx: context.Background(),
				input: auth.UserChangePinInput{
					UserPhoneNo:  validUserPhoneNo,
					OldPin:       validOldPin,
					NewPin:       validNewPin,
					RetypeNewPin: validRetypeNewPin,
					AuthUserId:   validUserId,
				},
			},
			expected: expected{
				Err:   auth.ErrPinCodeNotMatch,
				Value: auth.UserChangePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{Id: validUserPhoneNo, PinCode: pointer.String("asdfv")}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "failed user new pin and retype pin not match",
			args: args{
				ctx: context.Background(),
				input: auth.UserChangePinInput{
					UserPhoneNo:  validUserPhoneNo,
					OldPin:       validOldPin,
					NewPin:       validNewPin,
					RetypeNewPin: "asdfgh",
				},
			},
			expected: expected{
				Err:   auth.ErrPinCodeNotMatch,
				Value: auth.UserChangePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{Id: validUserPhoneNo, PinCode: &validOldPin}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}
			authService := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Logger{}}),
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
			})

			got, err := authService.UserChangePin(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
