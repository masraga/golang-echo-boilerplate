package auth

import (
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/dbtx"
)

type AuthRepository struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
}

type AuthRepositoryOpts struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
}

func NewAuthRepository(opts AuthRepositoryOpts) *AuthRepository {
	return &AuthRepository{
		DbTxInterface: opts.DbTxInterface,
		Sql:           opts.Sql,
	}
}
