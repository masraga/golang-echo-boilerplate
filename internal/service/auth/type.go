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
}

type CreateOTPInput struct {
	OtpCode       string
	UserId        string
	Note          *string
	ExpiredAtUtc0 int64
}

type CreateOTPOutput struct {
	OtpCode string
}

type FindOTPInput struct {
	UserId string
}

type FindOTPOutput struct {
	Id      string
	OtpCode string
}

type DeleteAllUserOTPInput struct {
	UserId string
}

type DeleteAllUserOTPOutput struct {
	IsSuccess bool
}
