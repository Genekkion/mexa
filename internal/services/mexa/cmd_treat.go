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

		casualty, err := s.repos.Casualties.GetCasualtyById(ctx, s.Exercise().Id, casualtyId)
		if err != nil {
			return err
		}

		return s.handleCasualtyCheck(ctx, u, casualty.FourD)
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

		var ikb [][]chatdomain.InlineKeyboardEntry

		const kbWidth = 2
		const kbHeight = 2
		for i := 0; i < kbHeight; i++ {
			kb := make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)

			for j := 0; j < kbWidth; j++ {
				p := i*kbWidth + j + 1
				kb = append(kb, chatdomain.InlineKeyboardEntry{
					Text:         fmt.Sprintf("P%d", p),
					CallbackData: fmt.Sprintf("%s::result:%d:%d", treatEndPrefix, casualtyId, p),
				})
			}

			ikb = append(ikb, kb)
		}

		return s.bot.Reply(ctx, u.ChatId(), "Please submit the outcome of treatment", chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
			InlineKeyboard: ikb,
		}))
	} else if strings.HasPrefix(data, "result:") {
		casualtyIdStr := strings.TrimPrefix(data, "result:")

		items := strings.Split(casualtyIdStr, ":")
		if len(items) != 2 {
			return fmt.Errorf("invalid treat end callback data: %s", data)
		}

		casualtyIdStr, outcomeStr := items[0], items[1]

		casualtyId, err := strconv.Atoi(casualtyIdStr)
		if err != nil {
			return err
		}

		pValue, err := strconv.Atoi(outcomeStr)
		if err != nil {
			return err
		}

		outcome, err := mexadomain.ParsePValue(pValue)
		if err != nil {
			return err
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
