package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	"strconv"
	"strings"
)

func (s *Service) callbackTreatStart(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, treatStartPrefix+"::")
	if strings.HasPrefix(data, "start:") {
		casualtyId, err := strconv.Atoi(strings.TrimPrefix(data, "start:"))
		if err != nil {
			return err
		}

		l := mexadomain.NewCCLogTreatStart(casualtyId)

		err = s.repos.CCLogs.AddLog(ctx, casualtyId, mexadomain.CCLogTypeTreatStart, l.Value)
		if err != nil {
			return err
		}

		err = s.bot.Reply(ctx, u.ChatId(), "Treatment started")
		if err != nil {
			return err
		}

		ikb, err := s.kbCasualtyCheckCasualty(ctx, u, casualtyId)
		if err != nil {
			return err
		}

		// TODO: GETTING 400 FROM BELOW
		return s.bot.EditMessage(ctx, u.ChatId(), u.CallbackQuery.Message.MessageId, u.CallbackQuery.Message.Text, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
			InlineKeyboard: ikb,
		}))
	}

	fmt.Printf("Unknown treat start callback: %s\n", u.CallbackQuery.Data)
	return nil
}

func (s *Service) callbackTreatEnd(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, treatEndPrefix+"::")
	if strings.HasPrefix(data, "end:") {
		casualtyId, err := strconv.Atoi(strings.TrimPrefix(data, "end:"))
		if err != nil {
			return err
		}

		ikb := [][]chatdomain.InlineKeyboardEntry{
			{
				{
					Text:         "Success",
					CallbackData: fmt.Sprintf("%s::result:%d:success", treatEndPrefix, casualtyId),
				},
				{
					Text:         "Failure",
					CallbackData: fmt.Sprintf("%s::result:%d:failure", treatEndPrefix, casualtyId),
				},
			},
		}

		return s.bot.Reply(ctx, u.ChatId(), "Please submit the outcome of treatment", chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
			InlineKeyboard: ikb,
		}))
	} else if strings.HasPrefix(data, "result:") {
		casualtyIdStr := strings.TrimPrefix(data, "result:")
		success := false
		casualtyIdStr, success = strings.CutSuffix(casualtyIdStr, ":success")
		casualtyIdStr = strings.TrimSuffix(casualtyIdStr, ":failure")

		casualtyId, err := strconv.Atoi(casualtyIdStr)
		if err != nil {
			return err
		}

		var outcome mexadomain.CCLogEndOutcome
		if success {
			outcome = mexadomain.CCLogEndOutcomeSuccess
		} else {
			outcome = mexadomain.CCLogEndOutcomeFailure
		}

		err = s.repos.CCLogs.AddLog(ctx, casualtyId, mexadomain.CCLogTypeTreatEnd, mexadomain.CCLogValue{
			Outcome: &outcome,
		})
		if err != nil {
			return err
		}

		return s.bot.Reply(ctx, u.ChatId(), "Treatment ended")
	}

	fmt.Printf("Unknown treat end callback: %s\n", u.CallbackQuery.Data)
	return nil
}
