package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

type CasesRepo interface {
	AddCase(ctx context.Context, exerciseId mexadomain.ExerciseId, value mexadomain.CaseValue) (id *mexadomain.CaseId, err error)
	GetCase(ctx context.Context, exerciseId mexadomain.ExerciseId, caseId mexadomain.CaseId) (res *mexadomain.Case, err error)
	GetCases(ctx context.Context, exerciseId mexadomain.ExerciseId) (res []mexadomain.Case, err error)
	ClearCases(ctx context.Context, exerciseId mexadomain.ExerciseId) (err error)
}
