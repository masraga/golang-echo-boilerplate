package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Connection struct {
	DBUrl string
}

func NewConnection(databaseUrl string) *Connection {
	return &Connection{
		DBUrl: databaseUrl,
	}
}

func (c *Connection) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
