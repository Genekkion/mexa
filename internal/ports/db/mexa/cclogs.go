package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

type CadetCaseLogsRepo interface {
	AddLog(ctx context.Context, casualtyId mexadomain.CasualtyId, logType mexadomain.CCLogType, logValue mexadomain.CCLogValue) (err error)
	GetLogsByCasualtyId(ctx context.Context, casualtyId mexadomain.CasualtyId) (res []mexadomain.CCLog, err error)
	GetLogsByExercise(ctx context.Context, exId mexadomain.ExerciseId) (res []mexadomain.CCLog, err error)
}
