package mexasqlite

import (
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func TestUsersRepo_CreateUserIfNotExists(t *testing.T) {
	db := sqlite.NewTestDb(t)
	defer db.Close()

	ctx := t.Context()
	baseRepo := sqlite.NewBaseRepo(db)
	repo := NewUsersRepo(&baseRepo)

	userId := 1234567890
	username := "testuser"

	// 1. Create user for the first time
	err := repo.CreateUserIfNotExists(ctx, userId, username)
	test.NilErr(t, err)

	// Verify user exists
	var count int
	err = db.Db().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", userId).Scan(&count)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 user only", 1, count)

	// 2. Try to create the same user again (should do nothing due to ON CONFLICT)
	err = repo.CreateUserIfNotExists(ctx, userId, "different_username")
	test.NilErr(t, err)

	// Verify still only 1 user exists and username hasn't changed (since it's DO NOTHING)
	err = db.Db().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", userId).Scan(&count)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 user only", 1, count)

	var dbUsername string
	err = db.Db().QueryRowContext(ctx, "SELECT username FROM users WHERE id = ?", userId).Scan(&dbUsername)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected username to be unchanged", username, dbUsername)
}
