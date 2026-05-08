package dbtx

import (
	"context"
	"database/sql"

	"github.com/leporo/sqlf"
)

type DbTx struct {
	Db *sql.DB
}

func (s *DbTx) Begin(ctx context.Context, opt *sql.TxOptions) (context.Context, error) {
	tx, err := s.Db.BeginTx(ctx, opt)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, DbTxId, tx), nil
}

func (s *DbTx) TxFromContext(ctx context.Context) (tx *sql.Tx, ok bool) {
	tx, ok = ctx.Value(DbTxId).(*sql.Tx)
	return tx, ok
}

func (s *DbTx) UseTx(ctx context.Context) sqlf.Executor {
	tx, ok := s.TxFromContext(ctx)
	if !ok {
		return s.Db
	}
	return tx
}

func (s *DbTx) CommitOrRollback(ctx context.Context, err error) error {
	tx, ok := s.TxFromContext(ctx)
	if !ok {
		return ErrFailedBeginTx
	}
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
