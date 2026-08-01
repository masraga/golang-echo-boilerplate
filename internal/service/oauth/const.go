package oauth

const (
	TableSchema    string = "public"
	OauthTableName string = TableSchema + "." + "oauth"
)

const (
	ExchangeTokenActionRegisterUser ExchangeTokenAction = "REGISTER_USER"
)

const (
	DefaultCreatedBy string = "oauth-system"
)
