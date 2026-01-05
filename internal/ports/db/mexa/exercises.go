package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

type ExercisesRepo interface {
	GetExerciseId(ctx context.Context, code string) (id *mexadomain.ExerciseId, err error)
	AddExercise(ctx context.Context, code string, name string) (id *mexadomain.ExerciseId, err error)
}
