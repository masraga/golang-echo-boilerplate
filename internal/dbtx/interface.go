package dbtx

import (
	"context"
	"database/sql"

	"github.com/leporo/sqlf"
)

type DbTxInterface interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (context.Context, error)
	TxFromContext(ctx context.Context) (tx *sql.Tx, ok bool)
	UseTx(ctx context.Context) sqlf.Executor
	CommitOrRollback(ctx context.Context, err error) error
}
