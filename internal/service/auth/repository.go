package auth

import (
	"database/sql"

	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
)

type AuthRepository struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
	Db  *sql.DB
	Err *ctxerr.CtxErr
}

type AuthRepositoryOpts struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
	Db  *sql.DB
	Err *ctxerr.CtxErr
}

func NewAuthRepository(opts AuthRepositoryOpts) *AuthRepository {
	return &AuthRepository{
		DbTxInterface: opts.DbTxInterface,
		Sql:           opts.Sql,
		Db:            opts.Db,
		Err:           opts.Err,
	}
}
