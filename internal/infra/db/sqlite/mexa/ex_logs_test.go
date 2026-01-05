package mexasqlite

import (
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/test"
	"testing"
)

func setupExLogsRepo(t *testing.T) (*sqlite.DB, *ExercisesRepo, *ExLogsRepo) {
	db := sqlite.NewTestDb(t)
	baseRepo := sqlite.NewBaseRepo(db)
	exRepo := NewExercisesRepo(&baseRepo)
	exLogsRepo := NewExLogsRepo(&baseRepo)
	return db, exRepo, exLogsRepo
}

func TestExerciseLogsRepo_AddAndGet(t *testing.T) {
	t.Parallel()
	db, exRepo, repo := setupExLogsRepo(t)
	defer db.Close()
	ctx := t.Context()

	exId, err := exRepo.AddExercise(ctx, "EX1", "Name")
	test.NilErr(t, err)

	userId := 12345
	err = repo.AddExLog(ctx, *exId, userId, mexadomain.LogTypeExStart)
	test.NilErr(t, err)

	logs, err := repo.GetAllExLogs(ctx, *exId)
	test.NilErr(t, err)
	test.AssertEqual(t, "expected 1 log", 1, len(logs))
	test.AssertEqual(t, "expected correct exerciseId", *exId, logs[0].ExerciseId)
	test.AssertEqual(t, "expected correct userId", userId, logs[0].UserId)
	test.AssertEqual(t, "expected correct type", mexadomain.LogTypeExStart, logs[0].Type)
}
