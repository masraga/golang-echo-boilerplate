package main

import (
	"database/sql"

	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/database"
	"github.com/masraga/kerp-api/internal/dbtx"
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
