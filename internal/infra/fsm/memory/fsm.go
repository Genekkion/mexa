package memory

import (
	chatdomain "mexa/internal/domains/chat"
	fsmports "mexa/internal/ports/fsm"
	"sync"
)

type Fsm struct {
	state fsmports.FsmState

	userFsm  map[chatdomain.UserId]*UserFsm
	userData map[chatdomain.UserId]any
	mu       *sync.RWMutex
}

func NewFsm() *Fsm {
	return &Fsm{
		state:    fsmports.StateExPreparing,
		userFsm:  make(map[chatdomain.UserId]*UserFsm),
		userData: make(map[chatdomain.UserId]any),
		mu:       &sync.RWMutex{},
	}
}

func (f *Fsm) UserFsm(userId chatdomain.UserId) (sm *UserFsm, ok bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	sm, ok = f.userFsm[userId]
	return sm, ok
}

func (f *Fsm) RegisterUser(userId chatdomain.UserId) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.userFsm[userId] = NewUserFsm()
}

func (f *Fsm) UnregisterUser(userId chatdomain.UserId) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.userFsm, userId)
}

func (f *Fsm) UserState(userId chatdomain.UserId) (state fsmports.UserState) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	sm, ok := f.userFsm[userId]
	if !ok {
		return fsmports.UserStateDefault
	}
	return sm.State()
}

func (f *Fsm) SetUserState(userId chatdomain.UserId, state fsmports.UserState) {
	f.mu.Lock()
	defer f.mu.Unlock()

	sm, ok := f.userFsm[userId]
	if !ok {
		sm = NewUserFsm()
		f.userFsm[userId] = sm
	}

	sm.SetState(state)
}

func (f *Fsm) UserData(userId chatdomain.UserId) (data any, ok bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	data, ok = f.userData[userId]

	return data, ok
}

func (f *Fsm) SetUserData(userId chatdomain.UserId, data any) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.userData[userId] = data
}

func (f *Fsm) DeleteUserData(userId chatdomain.UserId) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.userData, userId)
}

func (f *Fsm) FsmState() (res fsmports.FsmState) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	res = f.state

	return res
}

func (f *Fsm) SetState(s fsmports.FsmState) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.state = s
}

type UserFsm struct {
	mu    sync.RWMutex
	state fsmports.UserState
}

func NewUserFsm() *UserFsm {
	return &UserFsm{
		state: fsmports.UserStateDefault,
	}
}

func (fsm *UserFsm) State() fsmports.UserState {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	return fsm.state
}

func (fsm *UserFsm) SetState(s fsmports.UserState) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	fsm.state = s
}
