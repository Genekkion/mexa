package mexaservice

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	fsmports "mexa/internal/ports/fsm"
)

func (s *Service) cmdExStart(ctx context.Context, u chatdomain.Update) (err error) {
	err = s.repos.ExLogs.AddExLog(ctx, s.Exercise().Id, u.UserId(), mexadomain.LogTypeExStart)
	if err != nil {
		return err
	}

	s.fsm.SetState(fsmports.StateExStarted)

	return s.bot.Reply(ctx, u.ChatId(), "Exercise started")
}

func (s *Service) cmdExEnd(ctx context.Context, u chatdomain.Update) (err error) {
	err = s.repos.ExLogs.AddExLog(ctx, s.Exercise().Id, u.UserId(), mexadomain.LogTypeExEnd)
	if err != nil {
		return err
	}

	s.fsm.SetState(fsmports.StateExEnd)

	return s.bot.Reply(ctx, u.ChatId(), "Exercise ended")
}
