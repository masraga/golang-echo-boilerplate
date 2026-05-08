package auth

type JwtSecretType string

type JwtExpirationType int64

type CreateNewAccountInput struct {
	TokenType
	PhoneNo string
}

type CreateNewAccountOutput struct {
	TokenType
	Id string
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
