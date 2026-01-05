package sqlitemock

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
	mexaports "mexa/internal/ports/db/mexa"
)

type Repos struct {
	Transactional
	Users         UsersRepo
	Cases         CasesRepo
	Casualties    CasualtiesRepo
	Exercises     ExercisesRepo
	Deterioration DeteriorationRepo
	ExLogs        ExLogsRepo
	CCLogs        CCLogsRepo
}

type Transactional struct {
	CtxTxCalls      []context.Context
	TxRollbackCalls []context.Context
	TxCommitCalls   []context.Context
}

func (t *Transactional) CtxTx(ctx context.Context) (context.Context, error) {
	t.CtxTxCalls = append(t.CtxTxCalls, ctx)
	return ctx, nil
}

func (t *Transactional) TxRollback(ctx context.Context) error {
	t.TxRollbackCalls = append(t.TxRollbackCalls, ctx)
	return nil
}

func (t *Transactional) TxCommit(ctx context.Context) error {
	t.TxCommitCalls = append(t.TxCommitCalls, ctx)
	return nil
}

type UsersRepo struct {
	CreateUserIfNotExistsCalls []struct {
		Id       mexadomain.UserId
		Username string
	}
}

func (r *UsersRepo) CreateUserIfNotExists(ctx context.Context, id mexadomain.UserId, username string) error {
	r.CreateUserIfNotExistsCalls = append(r.CreateUserIfNotExistsCalls, struct {
		Id       mexadomain.UserId
		Username string
	}{id, username})
	return nil
}

type ExercisesRepo struct {
	GetExerciseIdFunc func(ctx context.Context, code string) (*mexadomain.ExerciseId, error)
	AddExerciseFunc   func(ctx context.Context, code string, name string) (*mexadomain.ExerciseId, error)
}

func (r *ExercisesRepo) GetExerciseId(ctx context.Context, code string) (*mexadomain.ExerciseId, error) {
	if r.GetExerciseIdFunc != nil {
		return r.GetExerciseIdFunc(ctx, code)
	}
	return nil, nil
}

func (r *ExercisesRepo) AddExercise(ctx context.Context, code string, name string) (*mexadomain.ExerciseId, error) {
	if r.AddExerciseFunc != nil {
		return r.AddExerciseFunc(ctx, code, name)
	}
	return nil, nil
}

type CasesRepo struct {
	AddCaseFunc    func(ctx context.Context, exerciseId mexadomain.ExerciseId, value mexadomain.CaseValue) (*mexadomain.CaseId, error)
	GetCaseFunc    func(ctx context.Context, exerciseId mexadomain.ExerciseId, caseId mexadomain.CaseId) (*mexadomain.Case, error)
	GetCasesFunc   func(ctx context.Context, exerciseId mexadomain.ExerciseId) ([]mexadomain.Case, error)
	ClearCasesFunc func(ctx context.Context, exerciseId mexadomain.ExerciseId) error
}

func (r *CasesRepo) AddCase(ctx context.Context, exerciseId mexadomain.ExerciseId, value mexadomain.CaseValue) (*mexadomain.CaseId, error) {
	if r.AddCaseFunc != nil {
		return r.AddCaseFunc(ctx, exerciseId, value)
	}
	return nil, nil
}

func (r *CasesRepo) GetCase(ctx context.Context, exerciseId mexadomain.ExerciseId, caseId mexadomain.CaseId) (*mexadomain.Case, error) {
	if r.GetCaseFunc != nil {
		return r.GetCaseFunc(ctx, exerciseId, caseId)
	}
	return nil, nil
}

func (r *CasesRepo) GetCases(ctx context.Context, exerciseId mexadomain.ExerciseId) ([]mexadomain.Case, error) {
	if r.GetCasesFunc != nil {
		return r.GetCasesFunc(ctx, exerciseId)
	}
	return nil, nil
}

func (r *CasesRepo) ClearCases(ctx context.Context, exerciseId mexadomain.ExerciseId) error {
	if r.ClearCasesFunc != nil {
		return r.ClearCasesFunc(ctx, exerciseId)
	}
	return nil
}

type CasualtiesRepo struct {
	AddCasualtyFunc       func(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D, caseId mexadomain.CaseId) (*mexadomain.CasualtyId, error)
	DeleteCasualtyFunc    func(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) error
	GetCasualtiesByExFunc func(ctx context.Context, exerciseId int) ([]mexadomain.Casualty, error)
	GetCasualtyByIdFunc   func(ctx context.Context, exerciseId int, casualtyId mexadomain.CasualtyId) (*mexadomain.Casualty, error)
	GetCasualtyBy4DFunc   func(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (*mexadomain.Casualty, error)
}

func (r *CasualtiesRepo) AddCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D, caseId mexadomain.CaseId) (*mexadomain.CasualtyId, error) {
	if r.AddCasualtyFunc != nil {
		return r.AddCasualtyFunc(ctx, exerciseId, cadet4D, caseId)
	}
	return nil, nil
}

func (r *CasualtiesRepo) DeleteCasualty(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) error {
	if r.DeleteCasualtyFunc != nil {
		return r.DeleteCasualtyFunc(ctx, exerciseId, cadet4D)
	}
	return nil
}

func (r *CasualtiesRepo) GetCasualtiesByEx(ctx context.Context, exerciseId int) ([]mexadomain.Casualty, error) {
	if r.GetCasualtiesByExFunc != nil {
		return r.GetCasualtiesByExFunc(ctx, exerciseId)
	}
	return nil, nil
}

func (r *CasualtiesRepo) GetCasualtyById(ctx context.Context, exerciseId int, casualtyId mexadomain.CasualtyId) (*mexadomain.Casualty, error) {
	if r.GetCasualtyByIdFunc != nil {
		return r.GetCasualtyByIdFunc(ctx, exerciseId, casualtyId)
	}
	return nil, nil
}

func (r *CasualtiesRepo) GetCasualtyBy4D(ctx context.Context, exerciseId int, cadet4D mexadomain.Cadet4D) (*mexadomain.Casualty, error) {
	if r.GetCasualtyBy4DFunc != nil {
		return r.GetCasualtyBy4DFunc(ctx, exerciseId, cadet4D)
	}
	return nil, nil
}

type DeteriorationRepo struct {
	AddDeteriorationFunc           func(ctx context.Context, casualtyId mexadomain.CasualtyId, value string) (*mexadomain.CaseDeteriorationId, error)
	GetDeteriorationByCasualtyFunc func(ctx context.Context, casualtyId mexadomain.CasualtyId) ([]mexadomain.CadetDeterioration, error)
}

func (r *DeteriorationRepo) AddDeterioration(ctx context.Context, casualtyId mexadomain.CasualtyId, value string) (*mexadomain.CaseDeteriorationId, error) {
	if r.AddDeteriorationFunc != nil {
		return r.AddDeteriorationFunc(ctx, casualtyId, value)
	}
	return nil, nil
}

func (r *DeteriorationRepo) GetDeteriorationByCasualty(ctx context.Context, casualtyId mexadomain.CasualtyId) ([]mexadomain.CadetDeterioration, error) {
	if r.GetDeteriorationByCasualtyFunc != nil {
		return r.GetDeteriorationByCasualtyFunc(ctx, casualtyId)
	}
	return nil, nil
}

type ExLogsRepo struct {
	AddExLogFunc     func(ctx context.Context, exerciseId mexadomain.ExerciseId, userId mexadomain.UserId, exType mexadomain.ExLogType) error
	GetAllExLogsFunc func(ctx context.Context, exerciseId mexadomain.ExerciseId) ([]mexadomain.ExLog, error)
}

func (r *ExLogsRepo) AddExLog(ctx context.Context, exerciseId mexadomain.ExerciseId, userId mexadomain.UserId, exType mexadomain.ExLogType) error {
	if r.AddExLogFunc != nil {
		return r.AddExLogFunc(ctx, exerciseId, userId, exType)
	}
	return nil
}

func (r *ExLogsRepo) GetAllExLogs(ctx context.Context, exerciseId mexadomain.ExerciseId) ([]mexadomain.ExLog, error) {
	if r.GetAllExLogsFunc != nil {
		return r.GetAllExLogsFunc(ctx, exerciseId)
	}
	return nil, nil
}

type CCLogsRepo struct {
	AddLogFunc              func(ctx context.Context, casualtyId mexadomain.CasualtyId, logType mexadomain.CCLogType, logValue mexadomain.CCLogValue) error
	GetLogsByCasualtyIdFunc func(ctx context.Context, casualtyId mexadomain.CasualtyId) ([]mexadomain.CCLog, error)
	GetLogsByExerciseFunc   func(ctx context.Context, exId mexadomain.ExerciseId) ([]mexadomain.CCLog, error)
}

func (r *CCLogsRepo) AddLog(ctx context.Context, casualtyId mexadomain.CasualtyId, logType mexadomain.CCLogType, logValue mexadomain.CCLogValue) error {
	if r.AddLogFunc != nil {
		return r.AddLogFunc(ctx, casualtyId, logType, logValue)
	}
	return nil
}

func (r *CCLogsRepo) GetLogsByCasualtyId(ctx context.Context, casualtyId mexadomain.CasualtyId) ([]mexadomain.CCLog, error) {
	if r.GetLogsByCasualtyIdFunc != nil {
		return r.GetLogsByCasualtyIdFunc(ctx, casualtyId)
	}
	return nil, nil
}

func (r *CCLogsRepo) GetLogsByExercise(ctx context.Context, exId mexadomain.ExerciseId) ([]mexadomain.CCLog, error) {
	if r.GetLogsByExerciseFunc != nil {
		return r.GetLogsByExerciseFunc(ctx, exId)
	}
	return nil, nil
}

var _ mexaports.UsersRepo = &UsersRepo{}
var _ mexaports.ExercisesRepo = &ExercisesRepo{}
var _ mexaports.CasesRepo = &CasesRepo{}
var _ mexaports.CasualtiesRepo = &CasualtiesRepo{}
var _ mexaports.CasualtyDeteriorationRepo = &DeteriorationRepo{}
var _ mexaports.ExLogsRepo = &ExLogsRepo{}
var _ mexaports.CadetCaseLogsRepo = &CCLogsRepo{}
