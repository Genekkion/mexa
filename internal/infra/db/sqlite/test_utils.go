package sqlite

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

func NewTestDb(t *testing.T) *DB {
	t.Helper()

	h := sha1.Sum([]byte(t.Name()))
	dsn := fmt.Sprintf("file:%x?mode=memory&cache=shared", h[:])

	db, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}

	db.db.SetMaxOpenConns(1)
	db.db.SetMaxIdleConns(1)

	ctx := t.Context()
	err = db.Init(ctx)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
