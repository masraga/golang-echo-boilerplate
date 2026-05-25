package main

import (
	"database/sql"
	"os"

	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/database"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/rs/zerolog"
)

func ProvideSqlDb(config *Config) (db *sql.DB, err error) {
	database := database.NewConnection(config.DatabaseUrl)
	db, err = database.Connect()
	return
}

func ProvideDbTx(db *sql.DB) dbtx.DbTxInterface {
	return &dbtx.DbTx{Db: db}
}

func ProvideSqlDialect() *sqlf.Dialect {
	return sqlf.PostgreSQL
}

func ProvideZerolog() zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}
