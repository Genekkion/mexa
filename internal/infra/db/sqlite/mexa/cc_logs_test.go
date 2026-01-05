package mexasqlite

import (
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupCCLogsRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo, *CasualtiesRepo, *CadetCaseLogsRepo) {
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	exRepo := NewExercisesRepo(&baseRepo)
	cadetsRepo := NewCasualtiesRepo(&baseRepo)
	ccLogsRepo := NewCadetCaseLogsRepo(&baseRepo)
	return db, exRepo, cadetsRepo, ccLogsRepo
}

func TestCadetCaseLogsRepo_AddAndGet(t *testing.T) {
	t.Parallel()
	db, exRepo, cadetsRepo, repo := setupCCLogsRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	casualtyId, err := cadetsRepo.AddCasualty(ctx, *exId, 123, 456)
	test.NilErr(t, err)

	err = repo.AddLog(ctx, *casualtyId, mexadomain.CCLogTypeTreatStart, mexadomain.CCLogValue{})
	test.NilErr(t, err)

	logs, err := repo.GetLogsByCasualtyId(ctx, *casualtyId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 log", 1, len(logs))
	test.AssertEqual(t, "expected correct casualtyId", *casualtyId, logs[0].CasualtyId)
	test.AssertEqual(t, "expected correct type", mexadomain.CCLogTypeTreatStart, logs[0].Type)
}
