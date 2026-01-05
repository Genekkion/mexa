package mexasqlite

import (
	"context"
	"encoding/json"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type CadetCaseLogsRepo struct {
	*sqlite.BaseRepo
}

func NewCadetCaseLogsRepo(base *sqlite.BaseRepo) (repo *CadetCaseLogsRepo) {
	return &CadetCaseLogsRepo{
		BaseRepo: base,
	}
}

func (repo *CadetCaseLogsRepo) AddLog(ctx context.Context, cadetId mexadomain.CasualtyId, logType mexadomain.CCLogType, logValue mexadomain.CCLogValue) (err error) {
	b, err := json.Marshal(logValue)
	if err != nil {
		return err
	}

	const stmt = "INSERT INTO casualties_case_logs (casualty_id, created_at, type, data) VALUES (?,?,?,?)"
	_, err = repo.Exec(ctx, stmt, cadetId, utils.TNow(), logType, b)
	return err
}

func (repo *CadetCaseLogsRepo) GetLogsByCasualtyId(ctx context.Context, cadetId mexadomain.CasualtyId) (res []mexadomain.CCLog, err error) {
	const stmt = "SELECT id, type, data FROM casualties_case_logs WHERE casualty_id=? ORDER BY created_at ASC"
	rows, err := repo.Query(ctx, stmt, cadetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := mexadomain.CCLog{
			CasualtyId: cadetId,
			Value:      mexadomain.CCLogValue{},
		}
		var data []byte
		err = rows.Scan(&l.Id, &l.Type, &data)
		if err != nil {
			return nil, err
		}

		if len(data) > 0 {
			err = json.Unmarshal(data, &l.Value)
			if err != nil {
				return nil, err
			}
		}

		res = append(res, l)
	}

	return res, nil
}

func (repo *CadetCaseLogsRepo) GetLogsByExercise(ctx context.Context, exId mexadomain.ExerciseId) (res []mexadomain.CCLog, err error) {
	const stmt = "SELECT casualties.Id, ccl.type, ccl.data FROM casualties_case_logs AS ccl INNER JOIN casualties ON casualties.id=casualty_id WHERE casualties.exercise_id=? ORDER BY created_at ASC"
	rows, err := repo.Query(ctx, stmt, exId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l mexadomain.CCLog
		var data []byte
		err = rows.Scan(&l.Id, &l.Type, &data)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &l.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, l)
	}

	return res, nil
}
