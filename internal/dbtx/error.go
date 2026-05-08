package dbtx

import "errors"

var (
	ErrFailedBeginTx = errors.New("Failed to begin db transaction")
)
