package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	fsmports "mexa/internal/ports/fsm"
	"strings"
)

func (s *Service) HandleText(ctx context.Context, u chatdomain.Update) (err error) {

	userState := s.fsm.UserState(u.UserId())

	switch userState {
	case fsmports.UserStateAttachingCase:
		return s.handleTextAttach(ctx, u)
	case fsmports.UserStateCheckingCasualty:
		return s.handleTextCasualtyCheck(ctx, u)
	case fsmports.UserStateAddDeteriorate:
		return s.handleTextAddDeterioration(ctx, u)
	}

	fmt.Printf("Unhandled text: %s\n", u.Message.Text)
	return nil
}

func (s *Service) HandleCallbacks(ctx context.Context, u chatdomain.Update) (err error) {
	for p, h := range s.callbacks {
		if strings.HasPrefix(u.CallbackQuery.Data, p) {
			return h(ctx, u)
		}
	}

	fmt.Printf("Unhandled callback: %s\n", u.CallbackQuery.Data)
	return nil
}

func (s *Service) HandleCommands(ctx context.Context, u chatdomain.Update) (err error) {
	err = s.repos.Users.CreateUserIfNotExists(ctx, u.UserId(), u.Username())
	if err != nil {
		return err
	}

	cmdStr := strings.TrimSpace(u.Message.Text)
	cmdStr = strings.TrimPrefix(cmdStr, "/")

	cmd, ok := s.commands[cmdStr]
	if !ok {
		return s.bot.Reply(ctx, u.ChatId(), "Unknown command")
	}

	return cmd.Handler(ctx, u)
}
