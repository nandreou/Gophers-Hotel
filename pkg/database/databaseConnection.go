package database

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func ConnectToDatabase(dbData string) (*DB, error) {
	conn, err := sql.Open("pgx", dbData)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	dbConn.SQL = conn

	return dbConn, nil
}
