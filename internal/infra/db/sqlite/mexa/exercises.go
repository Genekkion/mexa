package mexasqlite

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type ExercisesRepo struct {
	*sqlite.BaseRepo
}

func NewExercisesRepo(base *sqlite.BaseRepo) (repo *ExercisesRepo) {
	return &ExercisesRepo{
		BaseRepo: base,
	}
}

func (repo *ExercisesRepo) GetExerciseId(ctx context.Context, code string) (id *mexadomain.ExerciseId, err error) {
	stmt := "SELECT id FROM exercises WHERE code=?"
	err = repo.QueryRow(ctx, stmt, code).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (repo *ExercisesRepo) AddExercise(ctx context.Context, code string, name string) (id *mexadomain.ExerciseId, err error) {
	stmt := "INSERT INTO exercises (code, name, created_at) VALUES (?,?,?) ON CONFLICT DO UPDATE SET id=id RETURNING id"
	err = repo.QueryRow(ctx, stmt, code, name, utils.TNow()).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}
