package database_test

import (
	"testing"

	"github.com/masraga/golang-echo-boilerplate/internal/database"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Connection(t *testing.T) {
	dbObj := database.NewConnection(database.DATABASE_LOCAL_URL)
	_, err := dbObj.Connect()
	require.ErrorIs(t, err, nil)
}
