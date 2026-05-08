package auth

type TokenType string

// table definition
const (
	TableSchema string = "public."
	TableAuth   string = TableSchema + "auth"
)

// token type generator
const (
	TokenTypeJwt TokenType = "jwt"
)
