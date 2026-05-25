package auth

type TokenType string

// table definition
const (
	TableSchema       string = "public."
	TableAuth         string = TableSchema + "auth"
	AuthCodeTableName string = TableSchema + "auth_otp"
)

// credential config
const (
	TokenTypeJwt       TokenType = "jwt"
	OtpExpiredDuration int64     = 2 //expired in minutes
)
