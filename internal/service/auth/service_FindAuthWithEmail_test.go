package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_FindAuthWithEmail(t *testing.T) {
	var (
		expectedId      string = faker.UUIDHyphenated()
		expectedEmail   string = faker.Word()
		expectedPhoneNo string = faker.Word()
	)

	type args struct {
		ctx   context.Context
		input auth.FindAuthWithEmailInput
	}

	type fields struct {
		AuthRepositoryReader auth.AuthRepositoryReaderInterface
	}

	type expected = testutil.Result[auth.FindAuthWithEmailOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success find user by email",
			args: args{
				ctx: context.Background(),
				input: auth.FindAuthWithEmailInput{
					Email: expectedEmail,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.FindAuthWithEmailOutput{
					Email:   expectedEmail,
					PhoneNo: expectedPhoneNo,
					Id:      expectedId,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuthWithEmail(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthWithEmailOutput{
						Id:      expectedId,
						Email:   expectedEmail,
						PhoneNo: expectedPhoneNo,
					}, nil)
				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			svc := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
			})

			got, err := svc.FindAuthWithEmail(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
