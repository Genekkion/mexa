package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/utils"
	"strconv"
	"strings"
)

func (s *Service) getCasesList(ctx context.Context) (cs []mexadomain.Case, str string, err error) {
	cs, err = s.repos.Cases.GetCases(ctx, s.Exercise().Id)
	if err != nil {
		return nil, "", err
	}

	listStrings := make([]string, 0, 2*len(cs))
	listStrings = append(listStrings, "*Cases:*")
	for _, c := range cs {
		listStrings = append(listStrings,
			utils.EscapeMd2(
				fmt.Sprintf("%d. %s", c.Id, c.Summary),
			),
		)
	}

	return cs, strings.Join(listStrings, "\n"), nil
}

func (s *Service) paginatorCaseList(offset int, csLen int, prefix string) (res []chatdomain.InlineKeyboardEntry) {
	if offset == 0 && csLen > 9 {
		return []chatdomain.InlineKeyboardEntry{
			{
				Text:         "-",
				CallbackData: fmt.Sprintf("%s::ignore", prefix),
			},
			{
				Text:         "Next",
				CallbackData: fmt.Sprintf("%s::offset:1", prefix),
			},
		}
	} else if csLen >= (offset+1)*9 {
		return []chatdomain.InlineKeyboardEntry{
			{
				Text:         "Previous",
				CallbackData: fmt.Sprintf("%s::offset:%d", prefix, offset-1),
			},
			{
				Text:         "Next",
				CallbackData: fmt.Sprintf("%s::offset:%d", prefix, offset+1),
			},
		}
	}
	return []chatdomain.InlineKeyboardEntry{
		{
			Text:         "Previous",
			CallbackData: fmt.Sprintf("%s::offset:%d", prefix, offset-1),
		},
		{
			Text:         "-",
			CallbackData: fmt.Sprintf("%s::ignore", prefix),
		},
	}
}

func (s *Service) cmdListCases(ctx context.Context, u chatdomain.Update) (err error) {
	cs, str, err := s.getCasesList(ctx)
	if err != nil {
		return err
	}

	const kbWidth = 3
	const kbHeight = 3
	const kbArea = kbWidth * kbHeight

	ikb := make([][]chatdomain.InlineKeyboardEntry, 0, (len(cs)/kbHeight)+1)
	kb := make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
	for i := range min(len(cs), kbArea) {
		c := cs[i]
		kb = append(kb, chatdomain.InlineKeyboardEntry{
			Text:         fmt.Sprintf("Case #%d", c.Id),
			CallbackData: fmt.Sprintf("%s::info:%d", listCasesPrefix, c.Id),
		})
		if len(kb) == kbWidth {
			ikb = append(ikb, kb)
			kb = make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
		}
	}
	if len(kb) > 0 {
		ikb = append(ikb, kb)
	}

	if len(cs) > kbArea {
		kb = []chatdomain.InlineKeyboardEntry{
			{
				Text:         "-",
				CallbackData: fmt.Sprintf("%s::ignore", listCasesPrefix),
			},
			{
				Text:         "Next",
				CallbackData: fmt.Sprintf("%s::offset:1", listCasesPrefix),
			},
		}
		ikb = append(ikb, kb)
	}

	return s.bot.Reply(ctx, u.ChatId(), str, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
		InlineKeyboard: ikb,
	}))
}

func (s *Service) callbackListCases(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, listCasesPrefix+"::")
	if data == "ignore" {
		return nil

	} else if strings.HasPrefix(data, "offset:") {
		offsetStr := strings.TrimPrefix(data, "offset:")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return err
		}

		cs, str, err := s.getCasesList(ctx)
		if err != nil {
			return err
		}

		const kbWidth = 3
		const kbHeight = 3
		const kbArea = kbWidth * kbHeight

		ikb := make([][]chatdomain.InlineKeyboardEntry, 0, (len(cs)/kbHeight)+1)
		kb := make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
		for i := offset * kbArea; i < min(len(cs), (offset+1)*kbArea); i++ {
			c := cs[i]
			kb = append(kb, chatdomain.InlineKeyboardEntry{
				Text:         fmt.Sprintf("Case #%d", c.Id),
				CallbackData: fmt.Sprintf("%s::info:%d", listCasesPrefix, c.Id),
			})

			if len(kb) == kbWidth {
				ikb = append(ikb, kb)
				kb = make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
			}
		}

		kb = s.paginatorCaseList(offset, len(cs), listCasesPrefix)
		ikb = append(ikb, kb)

		return s.bot.EditMessage(ctx, u.ChatId(), u.CallbackQuery.Message.MessageId, str, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
			InlineKeyboard: ikb,
		}))

	} else if strings.HasPrefix(data, "info:") {
		id, err := strconv.Atoi(strings.TrimPrefix(data, "info:"))
		if err != nil {
			return err
		}

		c, err := s.repos.Cases.GetCase(ctx, s.Exercise().Id, id)
		if err != nil {
			return err
		}

		return s.bot.EditMessage(ctx, u.ChatId(), u.CallbackQuery.Message.MessageId, c.TgMd2())
	}

	fmt.Printf("Unknown list case callback: %s\n", u.CallbackQuery.Data)
	return nil
}
