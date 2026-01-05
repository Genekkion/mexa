package mexasqlite

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type ExLogsRepo struct {
	*sqlite.BaseRepo
}

func NewExLogsRepo(base *sqlite.BaseRepo) (repo *ExLogsRepo) {
	return &ExLogsRepo{
		BaseRepo: base,
	}
}

func (repo *ExLogsRepo) AddExLog(ctx context.Context, exerciseId mexadomain.ExerciseId, userId mexadomain.UserId, exType mexadomain.ExLogType) (err error) {
	stmt := "INSERT INTO exercise_logs (exercise_id, user_id, type, created_at) VALUES (?,?,?,?)"
	_, err = repo.Exec(ctx, stmt, exerciseId, userId, exType, utils.TNow())
	return err
}

func (repo *ExLogsRepo) GetAllExLogs(ctx context.Context, exerciseId mexadomain.ExerciseId) (res []mexadomain.ExLog, err error) {
	stmt := "SELECT id, user_id, type, created_at FROM exercise_logs WHERE exercise_id=? ORDER BY created_at ASC"
	rows, err := repo.Query(ctx, stmt, exerciseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l mexadomain.ExLog
		l.ExerciseId = exerciseId
		err = rows.Scan(&l.Id, &l.UserId, &l.Type, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, l)
	}

	return res, nil
}
