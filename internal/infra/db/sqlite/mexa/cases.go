package mexasqlite

import (
	"context"
	"encoding/json"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type CasesRepo struct {
	*sqlite.BaseRepo
}

func NewCasesRepo(base *sqlite.BaseRepo) (repo *CasesRepo) {
	return &CasesRepo{
		BaseRepo: base,
	}
}

func (repo *CasesRepo) AddCase(ctx context.Context, exerciseId mexadomain.ExerciseId, value mexadomain.CaseValue) (id *mexadomain.CaseId, err error) {
	b, err := value.Json()
	if err != nil {
		return nil, err
	}

	stmt := "INSERT INTO cases (id, created_at, exercise_id, value) VALUES ((SELECT COALESCE(MAX(id), 0) + 1 FROM cases WHERE exercise_id=?),?,?,?) RETURNING id"
	err = repo.QueryRow(ctx, stmt, exerciseId, utils.TNow(), exerciseId, b).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (repo *CasesRepo) GetCase(ctx context.Context, exerciseId mexadomain.ExerciseId, caseId mexadomain.CaseId) (res *mexadomain.Case, err error) {
	var b []byte
	stmt := "SELECT value FROM cases WHERE exercise_id=? AND id=? LIMIT 1"
	err = repo.QueryRow(ctx, stmt, exerciseId, caseId).Scan(&b)
	if err != nil {
		return nil, err
	}

	res = &mexadomain.Case{
		Id: caseId,
	}
	err = json.Unmarshal(b, &res.CaseValue)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *CasesRepo) GetCases(ctx context.Context, exerciseId mexadomain.ExerciseId) (res []mexadomain.Case, err error) {
	stmt := "SELECT id, value FROM cases WHERE exercise_id=? ORDER BY id ASC"
	rows, err := repo.Query(ctx, stmt, exerciseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c mexadomain.Case
		var value string
		err = rows.Scan(&c.Id, &value)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(value), &c.CaseValue)
		if err != nil {
			return nil, err
		}

		res = append(res, c)
	}

	return res, nil
}

func (repo *CasesRepo) ClearCases(ctx context.Context, exerciseId mexadomain.ExerciseId) (err error) {
	stmt := "DELETE FROM cases WHERE exercise_id=?"
	_, err = repo.Exec(ctx, stmt, exerciseId)
	return err
}
