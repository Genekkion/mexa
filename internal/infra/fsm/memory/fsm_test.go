package memory

import (
	fsmports "mexa/internal/ports/fsm"
	"mexa/internal/test"
	"testing"
)

func TestFsm_State(t *testing.T) {
	fsm := NewFsm()
	test.AssertEqual(t, "initial state should be Preparing", fsmports.StateExPreparing, fsm.FsmState())

	fsm.SetState(fsmports.StateExStarted)
	test.AssertEqual(t, "state should be Started", fsmports.StateExStarted, fsm.FsmState())
}

func TestFsm_UserState(t *testing.T) {
	fsm := NewFsm()
	userId := 123

	test.AssertEqual(t, "initial user state should be Default", fsmports.UserStateDefault, fsm.UserState(userId))

	fsm.SetUserState(userId, fsmports.UserStateAttachingCase)
	test.AssertEqual(t, "user state should be AttachingCase", fsmports.UserStateAttachingCase, fsm.UserState(userId))

	fsm.UnregisterUser(userId)
	test.AssertEqual(t, "user state after unregister should be Default", fsmports.UserStateDefault, fsm.UserState(userId))
}

func TestFsm_UserData(t *testing.T) {
	fsm := NewFsm()
	userId := 123

	_, ok := fsm.UserData(userId)
	test.Assert(t, "should not have user data initially", !ok)

	data := "some data"
	fsm.SetUserData(userId, data)

	got, ok := fsm.UserData(userId)
	test.Assert(t, "should have user data", ok)
	test.AssertEqual(t, "user data should match", data, got.(string))

	fsm.DeleteUserData(userId)
	_, ok = fsm.UserData(userId)
	test.Assert(t, "should not have user data after delete", !ok)
}
