package mexasqlite

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type CasualtiesDeteriorationRepo struct {
	*sqlite.BaseRepo
}

func NewCasualtiesDeteriorationRepo(base *sqlite.BaseRepo) (repo *CasualtiesDeteriorationRepo) {
	return &CasualtiesDeteriorationRepo{
		BaseRepo: base,
	}
}

func (repo *CasualtiesDeteriorationRepo) AddDeterioration(ctx context.Context, cadetId mexadomain.CasualtyId, value string) (id *mexadomain.CaseDeteriorationId, err error) {
	tNow := utils.TNow()
	stmt := "INSERT INTO casualty_case_deterioration (casualty_id, value, created_at) VALUES (?,?,?) RETURNING id"
	err = repo.QueryRow(ctx, stmt, cadetId, value, tNow).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (repo *CasualtiesDeteriorationRepo) GetDeteriorationByCasualty(ctx context.Context, cadetId mexadomain.CasualtyId) (res []mexadomain.CadetDeterioration, err error) {
	stmt := "SELECT id, value FROM casualty_case_deterioration WHERE casualty_id=?"
	rows, err := repo.Query(ctx, stmt, cadetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := mexadomain.CadetDeterioration{
			CadetId: cadetId,
		}
		err = rows.Scan(&d.Id, &d.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, d)
	}

	return res, nil
}
