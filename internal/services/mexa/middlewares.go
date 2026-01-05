package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	fsmports "mexa/internal/ports/fsm"
)

func (s *Service) wrapAdmin(f chatdomain.Handler) chatdomain.Handler {
	return func(ctx context.Context, u chatdomain.Update) error {
		ok := s.admins.Contains(u.ChatId())
		if !ok {
			return s.bot.Reply(ctx, u.ChatId(), "Unauthorized.")
		}
		return f(ctx, u)
	}
}

func (s *Service) wrapExStarted(f chatdomain.Handler) chatdomain.Handler {
	return func(ctx context.Context, u chatdomain.Update) error {
		state := s.fsm.FsmState()
		switch state {
		case fsmports.StateExPreparing:
			return s.bot.Reply(ctx, u.ChatId(), "Unable to process command, exercise not yet started")
		case fsmports.StateExEnd:
			return s.bot.Reply(ctx, u.ChatId(), "Unable to process command, exercise ended")

		case fsmports.StateExStarted:
			return f(ctx, u)
		}

		panic(fmt.Sprintf("Invalid state: %d", state))
	}
}
