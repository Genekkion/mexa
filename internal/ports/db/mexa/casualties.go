package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

type CasualtiesRepo interface {
	AddCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D, caseId mexadomain.CaseId) (id *mexadomain.CasualtyId, err error)
	DeleteCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (err error)
	GetCasualtiesByEx(ctx context.Context, exerciseId int) (res []mexadomain.Casualty, err error)
	GetCasualtyById(ctx context.Context, exerciseId int, casualtyId mexadomain.CasualtyId) (res *mexadomain.Casualty, err error)
	GetCasualtyBy4D(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (res *mexadomain.Casualty, err error)
}

type CasualtyDeteriorationRepo interface {
	AddDeterioration(ctx context.Context, casualtyId mexadomain.CasualtyId, value string) (id *mexadomain.CaseDeteriorationId, err error)
	GetDeteriorationByCasualty(ctx context.Context, casualtyId mexadomain.CasualtyId) (res []mexadomain.CadetDeterioration, err error)
}
