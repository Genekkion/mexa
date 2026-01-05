package mexasqlite

import (
	"context"
	"database/sql"
	"errors"
	"mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type CasualtiesRepo struct {
	*sqlite.BaseRepo
}

func (repo *CasualtiesRepo) GetCasualtyById(ctx context.Context, exerciseId int, casualtyId mexadomain.CasualtyId) (res *mexadomain.Casualty, err error) {
	const stmt = "SELECT id, four_d, case_id FROM casualties WHERE exercise_id=? AND id=?"
	res = &mexadomain.Casualty{
		ExerciseId: exerciseId,
	}
	err = repo.QueryRow(ctx, stmt, exerciseId, casualtyId).Scan(&res.Id, &res.FourD, &res.CaseId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mexadomain.NoCasualtyFoundError
		}
		return nil, err
	}

	return res, nil
}

func NewCasualtiesRepo(base *sqlite.BaseRepo) (repo *CasualtiesRepo) {
	return &CasualtiesRepo{
		BaseRepo: base,
	}
}

func (repo *CasualtiesRepo) GetCasualtiesByEx(ctx context.Context, exerciseId int) (res []mexadomain.Casualty, err error) {
	stmt := "SELECT id, four_d, case_id FROM casualties WHERE exercise_id=?"
	rows, err := repo.Query(ctx, stmt, exerciseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := mexadomain.Casualty{
			ExerciseId: exerciseId,
		}
		err = rows.Scan(&c.Id, &c.FourD, &c.CaseId)
		if err != nil {
			return nil, err
		}

		res = append(res, c)
	}

	return res, nil
}

func (repo *CasualtiesRepo) AddCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D, caseId mexadomain.CaseId) (id *mexadomain.CasualtyId, err error) {
	tNow := utils.TNow()
	stmt := "INSERT INTO casualties (exercise_id, four_d, case_id, created_at) VALUES (?,?,?,?) RETURNING id"
	err = repo.QueryRow(ctx, stmt, exerciseId, cadet4D, caseId, tNow).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (repo *CasualtiesRepo) DeleteCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (err error) {
	stmt := "DELETE FROM casualties WHERE exercise_id=? AND four_d=?"
	_, err = repo.Exec(ctx, stmt, exerciseId, cadet4D)
	return err
}

func (repo *CasualtiesRepo) GetCasualtyBy4D(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (res *mexadomain.Casualty, err error) {
	stmt := "SELECT id, four_d, case_id FROM casualties WHERE exercise_id=? AND four_d=?"
	res = &mexadomain.Casualty{
		ExerciseId: exerciseId,
	}
	err = repo.QueryRow(ctx, stmt, exerciseId, cadet4D).Scan(&res.Id, &res.FourD, &res.CaseId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mexadomain.NoCasualtyFoundError
		}
		return nil, err
	}

	return res, nil
}
