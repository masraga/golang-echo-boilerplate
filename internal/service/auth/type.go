package auth

type JwtSecretType string

type JwtExpirationType int64

type CreateNewAccountInput struct {
	TokenType
	Id      string
	PhoneNo string
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
