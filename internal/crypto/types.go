package crypto

type ConfigCryptoKey string

type EncryptInput struct {
	PlainText string
}

type EncryptOutput struct {
	Result string
}

type DEncryptInput struct {
	HashCode string
}

type DecryptOutput struct {
	Result string
}
