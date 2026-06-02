package auth

import (
	"context"

	"github.com/masraga/kerp-api/internal/dbtx"
)

type AuthServiceInterface interface {
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
	CreateToken(ctx context.Context, input UserTokenClaimInput) (output UserTokenClaimOutput, err error)
	CreateJWTToken(ctx context.Context, input CreateJWTTokenInput) (output CreateJWTTokenOutput, err error)
	CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error)
	VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error)
	VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error)
	AuthValidatePin(ctx context.Context, input AuthValidatePinInput) (output AuthValidatePinOutput, err error)
	ValidateJwtToken(ctx context.Context, input ValidateJwtTokenInput) (output ValidateJwtTokenOutput, err error)

	CreateAuthApiContract(ctx context.Context, input CreateAuthApiContractInput) (output CreateAuthApiContractOutput, err error)
	GetAuthApiContract(ctx context.Context, input GetAuthApiContractInput) (output GetAuthApiContractOutput, err error)
	ListAuthApiContracts(ctx context.Context, input ListAuthApiContractsInput) (output ListAuthApiContractsOutput, err error)
	UpdateAuthApiContract(ctx context.Context, input UpdateAuthApiContractInput) (output UpdateAuthApiContractOutput, err error)
	DeleteAuthApiContract(ctx context.Context, input DeleteAuthApiContractInput) (output DeleteAuthApiContractOutput, err error)

	CreateAuthUserApiContract(ctx context.Context, input CreateAuthUserApiContractInput) (output CreateAuthUserApiContractOutput, err error)
	GetAuthUserApiContract(ctx context.Context, input GetAuthUserApiContractInput) (output GetAuthUserApiContractOutput, err error)
	ListAuthUserApiContracts(ctx context.Context, input ListAuthUserApiContractsInput) (output ListAuthUserApiContractsOutput, err error)
	UpdateAuthUserApiContract(ctx context.Context, input UpdateAuthUserApiContractInput) (output UpdateAuthUserApiContractOutput, err error)
	DeleteAuthUserApiContract(ctx context.Context, input DeleteAuthUserApiContractInput) (output DeleteAuthUserApiContractOutput, err error)

	ValidateUserApiContract(ctx context.Context, input ValidateUserApiContractInput) (output ValidateUserApiContractOutput, err error)
	BootstrapUserApiContracts(ctx context.Context, input BootstrapUserApiContractsInput) (output BootstrapUserApiContractsOutput, err error)

	CreateAuthRole(ctx context.Context, input CreateAuthRoleInput) (output CreateAuthRoleOutput, err error)
	GetAuthRole(ctx context.Context, input GetAuthRoleInput) (output GetAuthRoleOutput, err error)
	ListAuthRoles(ctx context.Context, input ListAuthRolesInput) (output ListAuthRolesOutput, err error)
	UpdateAuthRole(ctx context.Context, input UpdateAuthRoleInput) (output UpdateAuthRoleOutput, err error)
	DeleteAuthRole(ctx context.Context, input DeleteAuthRoleInput) (output DeleteAuthRoleOutput, err error)

	CreateAuthRoleContractApi(ctx context.Context, input CreateAuthRoleContractApiInput) (output CreateAuthRoleContractApiOutput, err error)
	ListAuthRoleContractApis(ctx context.Context, input ListAuthRoleContractApisInput) (output ListAuthRoleContractApisOutput, err error)
	DeleteAuthRoleContractApi(ctx context.Context, input DeleteAuthRoleContractApiInput) (output DeleteAuthRoleContractApiOutput, err error)

	AssignAuthUserRole(ctx context.Context, input AssignAuthUserRoleInput) (output AssignAuthUserRoleOutput, err error)
	DeleteAuthUserRole(ctx context.Context, input DeleteAuthUserRoleInput) (output DeleteAuthUserRoleOutput, err error)
}

type AuthRepositoryWriterInterface interface {
	dbtx.DbTxInterface
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
	DeleteAllUserOTP(ctx context.Context, input DeleteAllUserOTPInput) (output DeleteAllUserOTPOutput, err error)
	CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error)
	VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error)
	VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error)
	CreateNewPin(ctx context.Context, input CreateNewPinInput) (output CreateNewPinOutput, err error)
	StoreAccessToken(ctx context.Context, input StoreAccessTokenInput) (output StoreAccessTokenOutput, err error)

	CreateAuthApiContract(ctx context.Context, input CreateAuthApiContractOutput) (output CreateAuthApiContractOutput, err error)
	UpdateAuthApiContract(ctx context.Context, input UpdateAuthApiContractInput) (output UpdateAuthApiContractOutput, err error)
	DeleteAuthApiContract(ctx context.Context, input DeleteAuthApiContractInput) (output DeleteAuthApiContractOutput, err error)

	CreateAuthUserApiContract(ctx context.Context, input CreateAuthUserApiContractOutput) (output CreateAuthUserApiContractOutput, err error)
	UpdateAuthUserApiContract(ctx context.Context, input UpdateAuthUserApiContractInput) (output UpdateAuthUserApiContractOutput, err error)
	DeleteAuthUserApiContract(ctx context.Context, input DeleteAuthUserApiContractInput) (output DeleteAuthUserApiContractOutput, err error)
	BootstrapUserApiContracts(ctx context.Context, input BootstrapUserApiContractsInput) (output BootstrapUserApiContractsOutput, err error)

	CreateAuthRole(ctx context.Context, input CreateAuthRoleOutput) (output CreateAuthRoleOutput, err error)
	UpdateAuthRole(ctx context.Context, input UpdateAuthRoleInput) (output UpdateAuthRoleOutput, err error)
	DeleteAuthRole(ctx context.Context, input DeleteAuthRoleInput) (output DeleteAuthRoleOutput, err error)

	CreateAuthRoleContractApi(ctx context.Context, input CreateAuthRoleContractApiOutput) (output CreateAuthRoleContractApiOutput, err error)
	DeleteAuthRoleContractApi(ctx context.Context, input DeleteAuthRoleContractApiInput) (output DeleteAuthRoleContractApiOutput, err error)

	DeleteAuthUserApiContractsByUserId(ctx context.Context, input DeleteAuthUserApiContractsByUserIdInput) (output DeleteAuthUserApiContractsByUserIdOutput, err error)
	InsertAuthUserApiContractsFromRole(ctx context.Context, input InsertAuthUserApiContractsFromRoleInput) (output InsertAuthUserApiContractsFromRoleOutput, err error)
	AssignAuthUserRole(ctx context.Context, input AssignAuthUserRoleInput) (output AssignAuthUserRoleOutput, err error)
	DeleteAuthUserRole(ctx context.Context, input DeleteAuthUserRoleInput) (output DeleteAuthUserRoleOutput, err error)
}

type AuthRepositoryReaderInterface interface {
	FindAuth(ctx context.Context, input FindAuthInput) (output FindAuthOutput, err error)
	FindOTP(ctx context.Context, input FindOTPInput) (output FindOTPOutput, err error)
	FindAccessToken(ctx context.Context, input FindAccessTokenInput) (output FindAccessTokenOutput, err error)

	GetAuthApiContract(ctx context.Context, input GetAuthApiContractInput) (output GetAuthApiContractOutput, err error)
	ListAuthApiContracts(ctx context.Context, input ListAuthApiContractsInput) (output ListAuthApiContractsOutput, err error)

	GetAuthUserApiContract(ctx context.Context, input GetAuthUserApiContractInput) (output GetAuthUserApiContractOutput, err error)
	ListAuthUserApiContracts(ctx context.Context, input ListAuthUserApiContractsInput) (output ListAuthUserApiContractsOutput, err error)
	ValidateUserApiContract(ctx context.Context, input ValidateUserApiContractInput) (output ValidateUserApiContractOutput, err error)

	GetAuthRole(ctx context.Context, input GetAuthRoleInput) (output GetAuthRoleOutput, err error)
	ListAuthRoles(ctx context.Context, input ListAuthRolesInput) (output ListAuthRolesOutput, err error)
	ListAuthRoleContractApis(ctx context.Context, input ListAuthRoleContractApisInput) (output ListAuthRoleContractApisOutput, err error)
}
