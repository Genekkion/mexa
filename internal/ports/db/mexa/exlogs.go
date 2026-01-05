package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

type ExLogsRepo interface {
	AddExLog(ctx context.Context, exerciseId mexadomain.ExerciseId, userId mexadomain.UserId, exType mexadomain.ExLogType) (err error)
	GetAllExLogs(ctx context.Context, exerciseId mexadomain.ExerciseId) (res []mexadomain.ExLog, err error)
}
