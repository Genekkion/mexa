package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	fsmports "mexa/internal/ports/fsm"
	"regexp"
	"strconv"
	"strings"
)

const (
	attachCaseMsg = "Attaching case for cadet:"
)

var (
	attachCaseMsgRegex = regexp.MustCompile(fmt.Sprintf(
		`%s\s+(\d{4})`,
		regexp.QuoteMeta(attachCaseMsg),
	))
)

func (s *Service) cmdAttach(ctx context.Context, u chatdomain.Update) (err error) {
	s.fsm.SetUserState(u.UserId(), fsmports.UserStateAttachingCase)

	return s.bot.Reply(ctx, u.ChatId(), "Enter Cadet's 4D number")
}

func (s *Service) handleTextAttach(ctx context.Context, u chatdomain.Update) error {
	cadet4dStr := strings.TrimSpace(u.Message.Text)
	matched := fourDRegex.MatchString(cadet4dStr)
	if !matched {
		return s.bot.Reply(ctx, u.ChatId(), "Invalid 4D number, please enter again")
	}

	cs, str, err := s.getCasesList(ctx)
	if err != nil {
		return err
	}

	str = fmt.Sprintf("%s %s\n\n%s", attachCaseMsg, cadet4dStr, str)

	const kbWidth = 3
	const kbHeight = 3
	const kbArea = kbWidth * kbHeight

	ikb := make([][]chatdomain.InlineKeyboardEntry, 0, (len(cs)/kbHeight)+1)
	kb := make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
	for i := range min(len(cs), kbArea) {
		c := cs[i]
		kb = append(kb, chatdomain.InlineKeyboardEntry{
			Text:         fmt.Sprintf("Case #%d", c.Id),
			CallbackData: fmt.Sprintf("%s::select:%d", attachCasePrefix, c.Id),
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
				CallbackData: ignorePrefix,
			},
			{
				Text:         "Next",
				CallbackData: fmt.Sprintf("%s::offset:1", attachCasePrefix),
			},
		}
		ikb = append(ikb, kb)
	}

	return s.bot.Reply(ctx, u.ChatId(), str, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
		InlineKeyboard: ikb,
	}))
}

func (s *Service) callbackAttachCase(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, attachCasePrefix+"::")
	if strings.HasPrefix(data, "select:") {
		return s.callbackAttachCaseSelect(ctx, u, data)
	} else if strings.HasPrefix(data, "offset:") {
		return s.callbackAttachCaseOffset(ctx, u, data)
	}

	return nil
}

func (s *Service) callbackAttachCaseOffset(ctx context.Context, u chatdomain.Update, data string) (err error) {
	data = strings.TrimPrefix(data, "offset:")
	offset, err := strconv.Atoi(data)
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
			CallbackData: fmt.Sprintf("%s::select:%d", attachCasePrefix, c.Id),
		})
		if len(kb) == kbWidth {
			ikb = append(ikb, kb)
			kb = make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
		}
	}

	kb = s.paginatorCaseList(offset, len(cs), attachCasePrefix)
	ikb = append(ikb, kb)

	return s.bot.EditMessage(ctx, u.ChatId(), u.CallbackQuery.Message.MessageId, str, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
		InlineKeyboard: ikb,
	}))
}

func (s *Service) callbackAttachCaseSelect(ctx context.Context, u chatdomain.Update, data string) (err error) {
	caseId, err := strconv.Atoi(strings.TrimPrefix(data, "select:"))
	if err != nil {
		return err
	}

	message4dStr := attachCaseMsgRegex.FindAllStringSubmatch(u.CallbackQuery.Message.Text, 1)
	if len(message4dStr) == 0 {
		fmt.Println("No 4Ds found in message")
		return s.bot.Reply(ctx, u.ChatId(), "Something went wrong, please try again")
	}

	cadet4d, err := strconv.Atoi(message4dStr[0][1])
	if err != nil {
		return err
	}

	_, err = s.repos.Casualties.GetCasualtyBy4D(ctx, s.Exercise().Id, cadet4d)
	if err == nil {
		s.fsm.SetUserState(u.UserId(), fsmports.UserStateDefault)
		return s.bot.Reply(ctx, u.ChatId(), "Unable to attach, this cadet already has a case attached")
	}

	_, err = s.repos.Casualties.AddCasualty(ctx, s.Exercise().Id, cadet4d, caseId)
	if err != nil {
		return err
	}

	s.fsm.SetUserState(u.UserId(), fsmports.UserStateDefault)

	c, err := s.repos.Cases.GetCase(ctx, s.Exercise().Id, caseId)
	if err != nil {
		return err
	}

	kb, err := s.kbCasualtyCheckCasualty(ctx, u, cadet4d)
	if err != nil {
		return err
	}

	return s.bot.Reply(ctx, u.ChatId(), c.TgMd2(), chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
		InlineKeyboard: kb,
	}))
}
