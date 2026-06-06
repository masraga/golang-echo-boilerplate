package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/external/fcm"
	"github.com/masraga/golang-echo-boilerplate/internal/database"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/service/notification"
	"github.com/rs/zerolog"
)

const errorLogFilePath = "internal/log/error.log"

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

func ProvideAuthAccessBootstrapUserId(configValue string) auth.AuthAccessBootstrapUserIdType {
	return auth.AuthAccessBootstrapUserIdType(configValue)
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

	if err := os.MkdirAll(filepath.Dir(errorLogFilePath), 0755); err != nil {
		return showLogMode
	}

	logFile, err := os.OpenFile(
		errorLogFilePath,
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

func ProvidePushNotificationService(ctx context.Context, cfg *Config) notification.PushProviderInterface {
	provider := fcm.NewFcmService(fcm.FcmServiceStructOpts{
		Ctx:              ctx,
		ServiceAccountId: cfg.FcmServiceAccountId,
	})
	return provider
}
