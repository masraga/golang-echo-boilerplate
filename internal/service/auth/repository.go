package auth

import (
	"database/sql"

	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/dbtx"
)

type AuthRepository struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
	Db  *sql.DB
}

type AuthRepositoryOpts struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
	Db  *sql.DB
}

func NewAuthRepository(opts AuthRepositoryOpts) *AuthRepository {
	return &AuthRepository{
		DbTxInterface: opts.DbTxInterface,
		Sql:           opts.Sql,
		Db:            opts.Db,
	}
}
