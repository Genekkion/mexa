package mexaservice

import (
	"context"
	"errors"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	fsmports "mexa/internal/ports/fsm"
	"mexa/internal/utils"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	fourDRegex = regexp.MustCompile(`^\d{4}$`)
)

func (s *Service) caseCheckListData(cs []mexadomain.Casualty) (res *string, rmk *chatdomain.ReplyMarkup, err error) {
	strs := make([]string, 0, len(cs)+2)
	strs = append(strs,
		"*Attached cases:*",
		"",
		"cadet\\_4d : case\\_id : status",
		"",
	)

	slices.SortFunc(cs, func(i mexadomain.Casualty, j mexadomain.Casualty) int {
		return i.FourD - j.FourD
	})

	for _, c := range cs {
		logs, err := s.repos.CCLogs.GetLogsByCasualtyId(context.Background(), c.Id)
		if err != nil {
			return nil, nil, err
		}

		var status string
		if len(logs) == 0 {
			status = "At frontline"
		} else {
			last := logs[len(logs)-1]
			switch last.Type {
			case mexadomain.CCLogTypeTreatStart:
				status = "At BCS"
			case mexadomain.CCLogTypeTreatEnd:
				status = last.Value.Outcome.String()
			default:
				return nil, nil, fmt.Errorf("unknown log type: %s", last.Type)
			}
		}
		status = fmt.Sprintf("*%s*", utils.EscapeMd2(status))

		strs = append(strs,
			fmt.Sprintf("*%d* \\: %d \\: %s", c.FourD, c.CaseId, status),
		)
	}
	strs = append(strs,
		"",
		"Select a cadet's 4D number below to open their case",
	)

	const kbWidth = 3
	const kbHeight = 3

	//const kbArea = kbWidth * kbHeight

	ikb := make([][]chatdomain.InlineKeyboardEntry, 0, (len(cs)/kbHeight)+1)
	kb := make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
	for _, c := range cs {
		kb = append(kb, chatdomain.InlineKeyboardEntry{
			Text: fmt.Sprintf("%d", c.FourD),
			CallbackData: fmt.Sprintf("%s::open:%d", casualtyCheckPrefix,
				c.FourD,
			),
		})
		if len(kb) == kbWidth {
			ikb = append(ikb, kb)
			kb = make([]chatdomain.InlineKeyboardEntry, 0, kbWidth)
		}
	}
	if len(kb) > 0 {
		ikb = append(ikb, kb)
	}

	str := strings.Join(strs, "\n")
	return &str, &chatdomain.ReplyMarkup{
		InlineKeyboard: ikb,
	}, nil
}

func (s *Service) cmdCasualties(ctx context.Context, u chatdomain.Update) (err error) {
	cs, err := s.repos.Casualties.GetCasualtiesByEx(ctx, s.Exercise().Id)
	if err != nil {
		return err
	} else if len(cs) == 0 {
		return s.bot.Reply(ctx, u.ChatId(), "No casualties found")
	}

	text, rmk, err := s.caseCheckListData(cs)
	if err != nil {
		return err
	}

	return s.bot.Reply(ctx, u.ChatId(), *text, chatdomain.WithReplyMarkup(*rmk))
}

func (s *Service) cmdCasualtyCheck(ctx context.Context, u chatdomain.Update) (err error) {
	s.fsm.SetUserState(u.UserId(), fsmports.UserStateCheckingCasualty)
	return s.bot.Reply(ctx, u.ChatId(), "Enter cadet's 4D number")
}

func (s *Service) handleTextCasualtyCheck(ctx context.Context, u chatdomain.Update) (err error) {
	cadet4dStr := strings.TrimSpace(u.Message.Text)
	matched := fourDRegex.MatchString(cadet4dStr)

	if !matched {
		return s.bot.Reply(ctx, u.ChatId(), "Invalid 4D number, please enter again")
	}

	cadet4d, err := strconv.Atoi(cadet4dStr)
	if err != nil {
		return err
	}

	s.fsm.SetUserState(u.UserId(), fsmports.UserStateDefault)

	casualty, err := s.repos.Casualties.GetCasualtyBy4D(ctx, s.Exercise().Id, cadet4d)
	if err != nil {
		if errors.Is(err, mexadomain.NoCasualtyFoundError) {
			return s.bot.Reply(ctx, u.ChatId(), "No casualty found for this 4D number")
		}

		return err
	}

	text := fmt.Sprintf("Case %d attached to cadet", casualty.CaseId)
	ikb := [][]chatdomain.InlineKeyboardEntry{
		{
			{
				Text:         "Open case",
				CallbackData: fmt.Sprintf("%s::open:%d", casualtyCheckPrefix, casualty.FourD),
			},
		},
	}

	return s.bot.Reply(ctx, u.ChatId(), text, chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
		InlineKeyboard: ikb,
	}))
}

func (s *Service) kbCasualtyCheckCasualty(ctx context.Context, _ chatdomain.Update, casualtyId mexadomain.CasualtyId) (res [][]chatdomain.InlineKeyboardEntry, err error) {
	logs, err := s.repos.CCLogs.GetLogsByCasualtyId(ctx, casualtyId)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		res = [][]chatdomain.InlineKeyboardEntry{
			{
				{
					Text:         "Start treatment timer",
					CallbackData: fmt.Sprintf("%s::start:%d", treatStartPrefix, casualtyId),
				},
			},
		}
	} else {
		res = [][]chatdomain.InlineKeyboardEntry{
			{
				{
					Text:         "Restart treatment timer",
					CallbackData: fmt.Sprintf("%s::start:%d", treatStartPrefix, casualtyId),
				},
				{
					Text:         "End treatment timer",
					CallbackData: fmt.Sprintf("%s::end:%d", treatEndPrefix, casualtyId),
				},
			},
		}
	}

	{

		deteriorationKb := []chatdomain.InlineKeyboardEntry{
			{
				Text:         "Add deterioration",
				CallbackData: fmt.Sprintf("%s::add:%d", deteriorationPrefix, casualtyId),
			},
		}

		res = append(res, deteriorationKb)
	}

	return res, nil
}

func (s *Service) callbackCasualtyCheck(ctx context.Context, u chatdomain.Update) (err error) {
	data := strings.TrimPrefix(u.CallbackQuery.Data, casualtyCheckPrefix+"::")
	if strings.HasPrefix(data, "open:") {
		cadet4d, err := strconv.Atoi(strings.TrimPrefix(data, "open:"))
		if err != nil {
			return err
		}

		c, err := s.repos.Casualties.GetCasualtyBy4D(ctx, s.Exercise().Id, cadet4d)
		if err != nil {
			return err
		}

		det, err := s.repos.Deterioration.GetDeteriorationByCasualty(ctx, c.Id)
		if err != nil {
			return err
		}

		cs, err := s.repos.Cases.GetCase(ctx, s.Exercise().Id, c.CaseId)
		if err != nil {
			return err
		}

		strs := []string{
			cs.TgMd2(),
			"__*Deterioration*__",
			"",
		}
		if len(det) > 0 {
			for _, d := range det {
				strs = append(strs, d.Value, "")
			}
		} else {
			strs = append(strs, "None")
		}

		ikb, err := s.kbCasualtyCheckCasualty(ctx, u, c.Id)
		if err != nil {
			return err
		}

		return s.bot.Reply(ctx, u.ChatId(), strings.Join(strs, "\n"), chatdomain.WithReplyMarkup(chatdomain.ReplyMarkup{
			InlineKeyboard: ikb,
		}))
	}

	fmt.Println("Unknown casualty check callback:", u.CallbackQuery.Data)
	return nil
}
