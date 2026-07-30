package auth

type TokenType string

// table definition
const (
	TableSchema       string = "public."
	TableAuth         string = TableSchema + "auth"
	AuthCodeTableName string = TableSchema + "auth_otp"
	AccessTokenTable  string = TableSchema + "auth_access_token"

	AuthApiContractTable      string = TableSchema + "auth_api_contract"
	AuthUserApiContractTable  string = TableSchema + "auth_user_api_contract"
	AuthRolesTable            string = TableSchema + "auth_roles"
	AuthRolesContractApiTable string = TableSchema + "auth_roles_contract_api"
)

// credential config
const (
	TokenTypeJwt            TokenType = "jwt"
	JwtTokenExpiredDuration int64     = 60 //expired in minutes
	OtpExpiredDuration      int64     = 2  //expired in minutes
	MaxPinLen               int       = 6
	MinPinLen               int       = 6
)

type LoginMode string

const (
	LOGIN_MODE_PHONE_NUMBER LoginMode = "PHONE_NUMBER"
	LOGIN_MODE_GOOGLE       LoginMode = "GOOGLE"
)

const (
	DEFAULT_PHONE_NO string = "0"
	DEFAULT_PIN_CODE string = "123456"
)
