package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	fsmports "mexa/internal/ports/fsm"
	"strconv"
	"strings"
)

func (s *Service) callbackDeterioration(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, deteriorationPrefix+"::")
	if strings.HasPrefix(data, "add:") {
		casualtyIdStr := strings.TrimPrefix(data, "add:")
		casualtyId, err := strconv.Atoi(casualtyIdStr)
		if err != nil {
			return err
		}

		defer func() {
			s.fsm.SetUserData(u.UserId(), casualtyId)
			s.fsm.SetUserState(u.UserId(), fsmports.UserStateAddDeteriorate)
		}()

		return s.bot.Reply(ctx, u.ChatId(), "Enter deterioration reason")
	}

	fmt.Printf("Unknown deterioration callback: %s\n", u.CallbackQuery.Data)
	return nil
}

func (s *Service) handleTextAddDeterioration(ctx context.Context, u chatdomain.Update) (err error) {
	s.fsm.SetUserState(u.UserId(), fsmports.UserStateDefault)
	casualtyIdAny, ok := s.fsm.UserData(u.UserId())
	if !ok {
		return s.bot.Reply(ctx, u.ChatId(), "Something went wrong, please try again")
	}
	defer s.fsm.DeleteUserData(u.UserId())

	casualtyId := casualtyIdAny.(mexadomain.CasualtyId)

	casualty, err := s.repos.Casualties.GetCasualtyById(ctx, s.Exercise().Id, casualtyId)
	if err != nil {
		return err
	}

	_, err = s.repos.Deterioration.AddDeterioration(ctx, casualtyId, strings.TrimSpace(u.Message.Text))
	if err != nil {
		return err
	}

	return s.handleCasualtyCheck(ctx, u, casualty.FourD)
}
