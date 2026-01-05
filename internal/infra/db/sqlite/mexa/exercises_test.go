package mexasqlite

import (
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupExercisesRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo) {
	t.Helper()
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	repo := NewExercisesRepo(&baseRepo)
	return db, repo
}

func TestExercisesRepo_AddExercise(t *testing.T) {
	t.Parallel()
	db, repo := setupExercisesRepo(t)
	defer db.Close()
	ctx := t.Context()

	code := "EX1"
	name := "Exercise 1"

	id, err := repo.AddExercise(ctx, code, name)
	test.NilErr(t, err)
	test.Assert(t, "expected id to be non-nil", id != nil)

	var dbName string
	err = db.Db().QueryRowContext(ctx, "SELECT name FROM exercises WHERE code=?", code).Scan(&dbName)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected correct name", name, dbName)
}

func TestExercisesRepo_AddExercise_Conflict(t *testing.T) {
	t.Parallel()
	db, repo := setupExercisesRepo(t)
	defer db.Close()
	ctx := t.Context()

	code := "EX1"
	name := "Exercise 1"

	id1, err := repo.AddExercise(ctx, code, name)
	test.NilErr(t, err)

	id2, err := repo.AddExercise(ctx, code, "different name")
	test.NilErr(t, err)
	test.AssertEqual(t, "expected same id from conflict", *id1, *id2)

	var dbName string
	err = db.Db().QueryRowContext(ctx, "SELECT name FROM exercises WHERE code=?", code).Scan(&dbName)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected name to be unchanged", name, dbName)
}

func TestExercisesRepo_GetExerciseId(t *testing.T) {
	t.Parallel()
	db, repo := setupExercisesRepo(t)
	defer db.Close()
	ctx := t.Context()

	code := "EX1"
	id, err := repo.AddExercise(ctx, code, "Name")
	test.NilErr(t, err)

	id2, err := repo.GetExerciseId(ctx, code)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected same id", *id, *id2)
}

func TestExercisesRepo_GetExerciseId_NotFound(t *testing.T) {
	t.Parallel()
	db, repo := setupExercisesRepo(t)
	defer db.Close()
	ctx := t.Context()

	id, err := repo.GetExerciseId(ctx, "NON_EXISTENT")
	test.Assert(t, "expected error for non-existent", err != nil)
	test.AssertEqual(t, "expected nil id for non-existent", (*mexadomain.ExerciseId)(nil), id)
}
