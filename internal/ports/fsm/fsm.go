package fsmports

import chatdomain "mexa/internal/domains/chat"

type FsmState int

const (
	StateExPreparing FsmState = iota
	StateExStarted
	StateExEnd
)

func (s FsmState) String() string {
	return []string{
		"Preparing",
		"Started",
		"Ended",
	}[s]
}

type Fsm interface {
	FsmState() (res FsmState)
	SetState(s FsmState)

	RegisterUser(userId chatdomain.UserId)
	UnregisterUser(userId chatdomain.UserId)
	UserState(userId chatdomain.UserId) (state UserState)
	SetUserState(userId chatdomain.UserId, state UserState)

	UserData(userId chatdomain.UserId) (data any, ok bool)
	SetUserData(userId chatdomain.UserId, data any)
	DeleteUserData(userId chatdomain.UserId)
}
