package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	botports "mexa/internal/ports/bot"
	fsmports "mexa/internal/ports/fsm"
	"mexa/internal/utils/set"
)

type Service struct {
	bot botports.Bot

	exercise mexadomain.Exercise
	batch    mexadomain.Batch

	commands  map[string]chatdomain.Command
	callbacks map[string]chatdomain.Handler
	admins    set.Set[chatdomain.UserId]
	fsm       fsmports.Fsm
	repos     Repos
}

type ServiceConfig struct {
	Bot      botports.Bot
	Exercise mexadomain.Exercise
	Batch    mexadomain.Batch
	Admins   []chatdomain.UserId
	Fsm      fsmports.Fsm
	Repos    Repos
}

func (c ServiceConfig) Validate() (err error) {
	if c.Bot == nil ||
		len(c.Admins) == 0 ||
		c.Fsm == nil {
		return fmt.Errorf("one or more required fields are nil")
	}
	err = c.Repos.Validate()
	if err != nil {
		return err
	}

	return nil
}

func NewService(ctx context.Context, c ServiceConfig) (ser *Service, err error) {
	err = c.Validate()
	if err != nil {
		return nil, err
	}

	ser = &Service{
		bot:       c.Bot,
		exercise:  c.Exercise,
		batch:     c.Batch,
		admins:    set.New(set.WithSlice(c.Admins)),
		commands:  make(map[string]chatdomain.Command),
		callbacks: make(map[string]chatdomain.Handler),
		fsm:       c.Fsm,
		repos:     c.Repos,
	}

	err = ser.initCommands(ctx)
	if err != nil {
		return nil, err
	}

	ser.initCallbacks()

	return ser, nil
}

func (s *Service) Exercise() mexadomain.Exercise {
	return s.exercise
}

func (s *Service) Batch() mexadomain.Batch {
	return s.batch
}
