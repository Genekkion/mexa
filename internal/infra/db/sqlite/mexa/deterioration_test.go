package mexasqlite

import (
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupDeteriorationRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo, *CasualtiesRepo, *CasualtiesDeteriorationRepo) {
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	exRepo := NewExercisesRepo(&baseRepo)
	cadetsRepo := NewCasualtiesRepo(&baseRepo)
	deteriorationRepo := NewCasualtiesDeteriorationRepo(&baseRepo)
	return db, exRepo, cadetsRepo, deteriorationRepo
}

func TestCadetsDeteriorationRepo_AddAndGet(t *testing.T) {
	t.Parallel()
	db, exRepo, casualtiesRepo, deteriorationRepo := setupDeteriorationRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	casualtyId, err := casualtiesRepo.AddCasualty(ctx, *exId, 123, 456)
	test.NilErr(t, err)

	val := "Condition worsened"
	detId, err := deteriorationRepo.AddDeterioration(ctx, *casualtyId, val)
	test.NilErr(t, err)
	test.Assert(t, "expected detId not nil", detId != nil)

	deteriorations, err := deteriorationRepo.GetDeteriorationByCasualty(ctx, *casualtyId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 deterioration", 1, len(deteriorations))
	test.AssertEqual(t, "expected correct value", val, deteriorations[0].Value)
	test.AssertEqual(t, "expected correct casualtyId", *casualtyId, deteriorations[0].CadetId)
}
