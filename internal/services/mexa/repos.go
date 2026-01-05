package mexaservice

import (
	"fmt"
	dbports "mexa/internal/ports/db"
	mexaports "mexa/internal/ports/db/mexa"
)

type Repos struct {
	dbports.Transactional
	Users         mexaports.UsersRepo
	Cases         mexaports.CasesRepo
	Casualties    mexaports.CasualtiesRepo
	Exercises     mexaports.ExercisesRepo
	Deterioration mexaports.CasualtyDeteriorationRepo
	ExLogs        mexaports.ExLogsRepo
	CCLogs        mexaports.CadetCaseLogsRepo
}

func (r *Repos) Validate() error {
	if r.Transactional == nil ||
		r.Users == nil ||
		r.Cases == nil ||
		r.Casualties == nil ||
		r.Exercises == nil ||
		r.Deterioration == nil ||
		r.ExLogs == nil ||
		r.CCLogs == nil {
		return fmt.Errorf("one or more repository fields are nil")
	}
	return nil
}
