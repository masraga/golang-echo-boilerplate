package main

import (
	"database/sql"
	"fmt"
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

func ProvideZerolog(config *Config) zerolog.Logger {
	showLogMode := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
	if config.ShowErrMode {
		return showLogMode
	}
	logFile, err := os.OpenFile(
		"internal/log/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		fmt.Println(err)
		return showLogMode
	}
	saveLogMode := zerolog.New(logFile).With().Timestamp().Caller().Logger()
	return saveLogMode
}
