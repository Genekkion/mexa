package fsmmock

import (
	chatdomain "mexa/internal/domains/chat"
	fsmports "mexa/internal/ports/fsm"
)

type Fsm struct {
	State       fsmports.FsmState
	UserStates  map[chatdomain.UserId]fsmports.UserState
	UserDataMap map[chatdomain.UserId]any
}

func New() *Fsm {
	return &Fsm{
		UserStates:  make(map[chatdomain.UserId]fsmports.UserState),
		UserDataMap: make(map[chatdomain.UserId]any),
	}
}

func (f *Fsm) FsmState() fsmports.FsmState {
	return f.State
}

func (f *Fsm) SetState(s fsmports.FsmState) {
	f.State = s
}

func (f *Fsm) RegisterUser(userId chatdomain.UserId) {
}

func (f *Fsm) UnregisterUser(userId chatdomain.UserId) {
	delete(f.UserStates, userId)
	delete(f.UserDataMap, userId)
}

func (f *Fsm) UserState(userId chatdomain.UserId) fsmports.UserState {
	return f.UserStates[userId]
}

func (f *Fsm) SetUserState(userId chatdomain.UserId, state fsmports.UserState) {
	f.UserStates[userId] = state
}

func (f *Fsm) UserData(userId chatdomain.UserId) (any, bool) {
	data, ok := f.UserDataMap[userId]
	return data, ok
}

func (f *Fsm) SetUserData(userId chatdomain.UserId, data any) {
	f.UserDataMap[userId] = data
}

func (f *Fsm) DeleteUserData(userId chatdomain.UserId) {
	delete(f.UserDataMap, userId)
}
