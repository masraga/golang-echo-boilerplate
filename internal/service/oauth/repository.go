package oauth

import (
	"database/sql"

	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
)

type OAuthRepository struct {
	dbtx.DbTxInterface
	sql *sqlf.Dialect
	db  *sql.DB
	err *ctxerr.CtxErr
}

type OAuthRepositoryOpts struct {
	dbtx.DbTxInterface
	Sql *sqlf.Dialect
	Db  *sql.DB
	Err *ctxerr.CtxErr
}

func NewOAuthRepository(opts OAuthRepositoryOpts) *OAuthRepository {
	return &OAuthRepository{
		DbTxInterface: opts.DbTxInterface,
		sql:           sqlf.PostgreSQL,
		db:            opts.Db,
		err:           opts.Err,
	}
}
