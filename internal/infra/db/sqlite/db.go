package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed sql/schema.sql
	schemaSql string
)

type DB struct {
	db *sql.DB
}

func New(fp string) (db *DB, err error) {
	db = &DB{}
	db.db, err = sql.Open("sqlite3", fp)
	if err != nil {
		return nil, err
	}

	err = db.db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Init(ctx context.Context) (err error) {
	fmt.Println("Initializing database")
	defer fmt.Println("Database initialized")

	_, err = db.db.ExecContext(ctx, schemaSql)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

// WARNING: Only to be used for testing purposes
func (db *DB) Db() *sql.DB {
	return db.db
}
