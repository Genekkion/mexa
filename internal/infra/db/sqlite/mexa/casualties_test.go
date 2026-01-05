package mexasqlite

import (
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupCasualtiesRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo, *CasualtiesRepo, *CasualtiesDeteriorationRepo) {
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	exRepo := NewExercisesRepo(&baseRepo)
	repo := NewCasualtiesRepo(&baseRepo)
	deteriorationRepo := NewCasualtiesDeteriorationRepo(&baseRepo)
	return db, exRepo, repo, deteriorationRepo
}

func TestCasualtiesRepo_AddAndGetCasualty(t *testing.T) {
	t.Parallel()
	db, exRepo, repo, _ := setupCasualtiesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	fourD := 123
	caseId := 456

	id, err := repo.AddCasualty(ctx, *exId, fourD, caseId)
	test.NilErr(t, err)
	test.Assert(t, "expected id not nil", id != nil)

	casualties, err := repo.GetCasualtiesByEx(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 casualty", 1, len(casualties))
	test.AssertEqual(t, "expected correct fourD", fourD, casualties[0].FourD)
	test.AssertEqual(t, "expected correct caseId", caseId, casualties[0].CaseId)
}

func TestCasualtiesRepo_DeleteCasualty(t *testing.T) {
	t.Parallel()
	db, exRepo, repo, _ := setupCasualtiesRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	fourD := 123
	_, err = repo.AddCasualty(ctx, *exId, fourD, 456)
	test.NilErr(t, err)

	err = repo.DeleteCasualty(ctx, *exId, fourD)
	test.NilErr(t, err)

	casualties, err := repo.GetCasualtiesByEx(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 0 casualties", 0, len(casualties))
}
