package fsmports

type UserState int

const (
	UserStateDefault UserState = iota
	UserStateCheckingCasualty
	UserStateAttachingCase
	UserStateAddDeteriorate
)
