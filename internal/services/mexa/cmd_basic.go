package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	"mexa/internal/utils"
	"strings"
)

func (s *Service) cmdQuit(ctx context.Context, u chatdomain.Update) (err error) {
	s.fsm.UnregisterUser(u.UserId())
	s.fsm.RegisterUser(u.UserId())

	return s.bot.Reply(ctx, u.ChatId(), "Cleared")
}

func (s *Service) cmdUser(ctx context.Context, u chatdomain.Update) (err error) {
	return s.bot.Reply(ctx, u.ChatId(), strings.Join([]string{
		"*ID*",
		fmt.Sprintf("%d", u.Message.From.Id),
	}, "\n"))
}

func (s *Service) cmdBatch(ctx context.Context, u chatdomain.Update) (err error) {
	batch := s.Batch()
	return s.bot.Reply(ctx, u.ChatId(), strings.Join([]string{
		"__*Batch info*__",
		"",
		"*Code*",
		utils.EscapeMd2(batch.Code),
		"",
		"*Name*",
		utils.EscapeMd2(batch.Name),
	}, "\n"))
}

func (s *Service) cmdExercise(ctx context.Context, u chatdomain.Update) (err error) {
	state := s.fsm.FsmState()
	exercise := s.Exercise()
	return s.bot.Reply(ctx, u.ChatId(), strings.Join([]string{
		"__*Exercise info*__",
		"",
		"*Code*",
		utils.EscapeMd2(exercise.Code),
		"",
		"*Name*",
		utils.EscapeMd2(exercise.Name),
		"",
		"*State*",
		fmt.Sprintf("%s", state.String()),
	}, "\n"))
}
