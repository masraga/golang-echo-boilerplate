package auth

type TokenType string

// table definition
const (
	TableSchema       string = "public."
	TableAuth         string = TableSchema + "auth"
	AuthCodeTableName string = TableSchema + "auth_otp"
	AccessTokenTable  string = TableSchema + "access_token"
)

// credential config
const (
	TokenTypeJwt            TokenType = "jwt"
	JwtTokenExpiredDuration int64     = 60 //expired in minutes
	OtpExpiredDuration      int64     = 2  //expired in minutes
	MaxPinLen               int       = 6
	MinPinLen               int       = 6
)
