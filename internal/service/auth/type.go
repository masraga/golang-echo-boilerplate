package auth

import "github.com/golang-jwt/jwt/v5"

type JwtSecretType string

type JwtExpirationType int64

type AuthAccessBootstrapUserIdType string

type CreateNewAccountInput struct {
	TokenType
	Id         string
	PhoneNo    string
	FirebaseId *string
}

type CreateNewAccountOutput struct {
	TokenType
	Id      string
	OtpCode string
}

type UserTokenClaimInput struct {
	TokenType
	UserId        string
	ExpiredAtUtc0 int64
	IssuerAtUtc0  int64
	UserName      string
}

type UserTokenClaimOutput struct {
	TokenType
	Token string
}

type CreateJWTTokenInput struct {
	ExpiredAtUtc0 int64
	IssuerAtUtc0  int64
	UserId        string
	Metadata      CreateJWTTokenMetadata
}

type CreateJWTTokenMetadata struct {
}

type CreateJWTTokenOutput struct {
	UserId string
	Token  string
}

type FindAuthInput struct {
	UserId  string
	PhoneNo string
}

type FindAuthOutput struct {
	Id      string
	PhoneNo string
	PinCode *string
}

type UpdateFirebaseIdInput struct {
	UserId     string
	FirebaseId string
}

type UpdateFirebaseIdOutput struct {
	UserId     string
	FirebaseId string
}

type CreateOTPInput struct {
	OtpCode       string
	PhoneNo       string
	UserId        string
	Note          *string
	ExpiredAtUtc0 int64
}

type CreateOTPOutput struct {
	OtpCode string
}

type FindOTPInput struct {
	UserId  string
	OtpCode string
}

type FindOTPOutput struct {
	Id            string
	OtpCode       string
	Note          *string
	ExpiredAtUtc0 int64
	IsVerified    bool
}

type DeleteAllUserOTPInput struct {
	UserId string
}

type DeleteAllUserOTPOutput struct {
	IsSuccess bool
}

type VerifyOtpInput struct {
	UserId  string
	PhoneNo string
	OtpCode string
}

type VerifyOtpOutput struct {
	IsValid bool
	UserId  string
	Note    *string
}

type VerifyUserAccountInput struct {
	PhoneNo string
}

type VerifyUserAccountOutput struct {
	UserId     string
	PhoneNo    string
	IsVerified bool
	IsNewUser  bool
}

type AuthValidatePinInput struct {
	UserId        string
	PhoneNo       string
	PinCode       string
	RetypePinCode *string
}

type AuthValidatePinOutput struct {
	IsValid       bool
	Token         string
	UserId        string
	ExpiredAtUtc0 int64
}

type CreateNewPinInput struct {
	PinCode string
	UserId  string
}

type CreateNewPinOutput struct {
	PinCode string
	UserId  string
}

type StoreAccessTokenInput struct {
	Token         string
	UserId        string
	ExpiredAtUtc0 int64
}

type StoreAccessTokenOutput struct {
	Token         string
	UserId        string
	ExpiredAtUtc0 int64
	IsActive      bool
}

type FindAccessTokenInput struct {
	Token  string
	UserId string
}

type FindAccessTokenOutput struct {
	Token         string
	UserId        string
	ExpiredAtUtc0 int64
	IsActive      bool
}

type ValidateJwtTokenInput struct {
	Token string
}

type ValidateJwtTokenOutput struct {
	jwt.RegisteredClaims

	UserId        string `json:"userId"`
	ExpiredAtUtc0 int64  `json:"exp"`
	IssuerAtUtc0  int64  `json:"iat"`
}

type AuthApiContract struct {
	Id             string
	EndpointPath   string
	EndpointMethod string
	Description    string
	CreatedAtUtc0  int64
	UpdatedAtUtc0  int64
	IsActive       bool
}

type CreateAuthApiContractInput struct {
	Id             string
	EndpointPath   string
	EndpointMethod string
	Description    string
}

type CreateAuthApiContractOutput = AuthApiContract

type GetAuthApiContractInput struct {
	Id string
}

type GetAuthApiContractOutput = AuthApiContract

type ListAuthApiContractsInput struct {
}

type ListAuthApiContractsOutput struct {
	Data []AuthApiContract
}

type UpdateAuthApiContractInput struct {
	Id             string
	EndpointPath   string
	EndpointMethod string
	Description    string
	IsActive       bool
}

type UpdateAuthApiContractOutput = AuthApiContract

type DeleteAuthApiContractInput struct {
	Id string
}

type DeleteAuthApiContractOutput struct {
	IsSuccess bool
}

type AuthUserApiContract struct {
	Id            string
	UserId        string
	ApiContractId string
	CreatedAtUtc0 int64
	UpdatedAtUtc0 int64
	IsActive      bool
}

type CreateAuthUserApiContractInput struct {
	UserId        string
	ApiContractId string
}

type CreateAuthUserApiContractOutput = AuthUserApiContract

type GetAuthUserApiContractInput struct {
	Id string
}

type GetAuthUserApiContractOutput = AuthUserApiContract

type ListAuthUserApiContractsInput struct {
}

type ListAuthUserApiContractsOutput struct {
	Data []AuthUserApiContract
}

type UpdateAuthUserApiContractInput struct {
	Id            string
	UserId        string
	ApiContractId string
	IsActive      bool
}

type UpdateAuthUserApiContractOutput = AuthUserApiContract

type DeleteAuthUserApiContractInput struct {
	Id string
}

type DeleteAuthUserApiContractOutput struct {
	IsSuccess bool
}

type DeleteAuthUserApiContractsByUserIdInput struct {
	UserId string
}

type DeleteAuthUserApiContractsByUserIdOutput struct {
	DeletedCount int64
}

type InsertAuthUserApiContractsFromRoleInput struct {
	UserId    string
	RoleId    string
	CreatedBy string
}

type InsertAuthUserApiContractsFromRoleOutput struct {
	InsertedCount int64
}

type ValidateUserApiContractInput struct {
	UserId         string
	EndpointPath   string
	EndpointMethod string
}

type ValidateUserApiContractOutput struct {
	IsAllowed bool
}

type BootstrapUserApiContractsInput struct {
	UserId string
}

type BootstrapUserApiContractsOutput struct {
	InsertedCount int64
}

type AuthRole struct {
	Id            string
	RoleName      string
	Description   string
	OwnerId       string
	CreatedAtUtc0 int64
	UpdatedAtUtc0 int64
	CreatedBy     string
	IsActive      bool
}

type CreateAuthRoleInput struct {
	RoleName    string
	Description string
	OwnerId     string
	CreatedBy   string
}

type CreateAuthRoleOutput = AuthRole

type GetAuthRoleInput struct {
	Id string
}

type GetAuthRoleOutput = AuthRole

type ListAuthRolesInput struct {
}

type ListAuthRolesOutput struct {
	Data []AuthRole
}

type UpdateAuthRoleInput struct {
	Id          string
	RoleName    string
	Description string
	OwnerId     string
	CreatedBy   string
	IsActive    bool
}

type UpdateAuthRoleOutput = AuthRole

type DeleteAuthRoleInput struct {
	Id string
}

type DeleteAuthRoleOutput struct {
	IsSuccess bool
}

type AuthRoleContractApi struct {
	Id                string
	RoleId            string
	AuthApiContractId string
	CreatedAtUtc0     int64
	UpdatedAtUtc0     int64
	CreatedBy         string
	IsActive          bool
}

type CreateAuthRoleContractApiInput struct {
	RoleId            string
	AuthApiContractId string
	CreatedBy         string
}

type CreateAuthRoleContractApiOutput = AuthRoleContractApi

type ListAuthRoleContractApisInput struct {
	RoleId string
}

type ListAuthRoleContractApisOutput struct {
	Data []AuthRoleContractApi
}

type DeleteAuthRoleContractApiInput struct {
	Id     string
	RoleId string
}

type DeleteAuthRoleContractApiOutput struct {
	IsSuccess bool
}

type AssignAuthUserRoleInput struct {
	UserId    string
	RoleId    string
	RoleName  string
	CreatedBy string
}

type AssignAuthUserRoleOutput struct {
	UserId        string
	RoleId        string
	RoleName      string
	GrantedCount  int64
	UpdatedAtUtc0 int64
	CreatedBy     string
}

type DeleteAuthUserRoleInput struct {
	UserId    string
	CreatedBy string
}

type DeleteAuthUserRoleOutput struct {
	IsSuccess bool
}
