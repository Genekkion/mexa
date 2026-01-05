package mexasqlite

import (
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupCasesRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo, *CasesRepo) {
	t.Helper()
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	exRepo := NewExercisesRepo(&baseRepo)
	repo := NewCasesRepo(&baseRepo)
	return db, exRepo, repo
}

func TestCasesRepo_AddAndGetCase(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupCasesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	val := mexadomain.CaseValue{
		Summary: "Summary",
	}
	id, err := repo.AddCase(ctx, *exId, val)
	test.NilErr(t, err)
	test.Assert(t, "expected id not nil", id != nil)

	c, err := repo.GetCase(ctx, *exId, *id)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected correct id", *id, c.Id)
	test.AssertEqual(t, "expected correct summary", val.Summary, c.CaseValue.Summary)
}

func TestCasesRepo_GetCases_Empty(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupCasesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	cases, err := repo.GetCases(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 0 cases", 0, len(cases))
}

func TestCasesRepo_GetCases_Ordering(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupCasesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	id1, err := repo.AddCase(ctx, *exId, mexadomain.CaseValue{Summary: "S1"})
	test.NilErr(t, err)
	id2, err := repo.AddCase(ctx, *exId, mexadomain.CaseValue{Summary: "S2"})
	test.NilErr(t, err)

	test.AssertEqual(t, "expected id2 > id1", true, *id2 > *id1)

	cases, err := repo.GetCases(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 2 cases", 2, len(cases))
	test.AssertEqual(t, "expected ASC order", true, cases[0].Id < cases[1].Id)
	test.AssertEqual(t, "first should be id1", *id1, cases[0].Id)
	test.AssertEqual(t, "second should be id2", *id2, cases[1].Id)
}

func TestCasesRepo_ClearCases(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupCasesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	_, err = repo.AddCase(ctx, *exId, mexadomain.CaseValue{
		Summary: "S",
	})
	test.NilErr(t, err)

	err = repo.ClearCases(ctx, *exId)
	test.NilErr(t, err)

	cases, err := repo.GetCases(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 0 cases", 0, len(cases))
}

func TestCasesRepo_GetCase_NotFound(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupCasesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	c, err := repo.GetCase(ctx, *exId, 999)
	test.Assert(t, "expected error", err != nil)
	test.AssertEqual(t, "expected nil case", (*mexadomain.Case)(nil), c)
}
