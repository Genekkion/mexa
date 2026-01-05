package mexaservice

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	botmock "mexa/internal/infra/bot/mock"
	sqlitemock "mexa/internal/infra/db/mock"
	fsmmock "mexa/internal/infra/fsm/mock"
	fsmports "mexa/internal/ports/fsm"
	"mexa/internal/test"
	"strconv"
	"testing"
)

func setupServiceTest(t *testing.T) (*Service, *botmock.Bot, *fsmmock.Fsm, *sqlitemock.Repos) {
	bot := &botmock.Bot{}
	fsm := fsmmock.New()
	repos := &sqlitemock.Repos{}

	config := ServiceConfig{
		Bot:    bot,
		Admins: []chatdomain.UserId{123},
		Fsm:    fsm,
		Repos: Repos{
			Transactional: &repos.Transactional,
			Users:         &repos.Users,
			Cases:         &repos.Cases,
			Casualties:    &repos.Casualties,
			Exercises:     &repos.Exercises,
			Deterioration: &repos.Deterioration,
			ExLogs:        &repos.ExLogs,
			CCLogs:        &repos.CCLogs,
		},
		Exercise: mexadomain.Exercise{Id: 1},
	}

	ctx := context.Background()
	svc, err := NewService(ctx, config)
	test.NilErr(t, err)

	return svc, bot, fsm, repos
}

func TestService_cmdExStart(t *testing.T) {
	svc, bot, fsm, repos := setupServiceTest(t)
	ctx := context.Background()

	// Mocking successful log addition
	repos.ExLogs.AddExLogFunc = func(ctx context.Context, exerciseId mexadomain.ExerciseId, userId mexadomain.UserId, exType mexadomain.ExLogType) error {
		test.AssertEqual(t, "exercise id should match", mexadomain.ExerciseId(1), exerciseId)
		test.AssertEqual(t, "user id should match", 123, userId)
		test.AssertEqual(t, "log type should be ExStart", mexadomain.LogTypeExStart, exType)
		return nil
	}

	u := chatdomain.Update{
		Message: &chatdomain.Message{
			From: chatdomain.User{Id: 123},
			Chat: struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			}{Id: 123},
			Text: "/ex_start",
		},
	}

	err := svc.HandleCommands(ctx, u)
	test.NilErr(t, err)

	test.AssertEqual(t, "FSM state should be Started", fsmports.StateExStarted, fsm.FsmState())
	test.AssertEqual(t, "should have 1 reply", 1, len(bot.ReplyCalls))
	test.AssertEqual(t, "reply text should match", "Exercise started", bot.ReplyCalls[0].Text)
}

func TestService_cmdExStart_Unauthorized(t *testing.T) {
	svc, bot, fsm, _ := setupServiceTest(t)
	ctx := context.Background()

	u := chatdomain.Update{
		Message: &chatdomain.Message{
			From: chatdomain.User{Id: 456},
			Chat: struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			}{Id: 456},
			Text: "/ex_start",
		},
	}

	err := svc.HandleCommands(ctx, u)
	test.NilErr(t, err)

	test.AssertEqual(t, "FSM state should still be Preparing", fsmports.StateExPreparing, fsm.FsmState())
	test.AssertEqual(t, "should have 1 reply", 1, len(bot.ReplyCalls))
	test.AssertEqual(t, "reply text should match", "Unauthorized.", bot.ReplyCalls[0].Text)
}

func TestService_callbackTreatStart(t *testing.T) {
	svc, bot, _, repos := setupServiceTest(t)
	ctx := context.Background()

	casualtyId := 789
	repos.CCLogs.AddLogFunc = func(ctx context.Context, id mexadomain.CasualtyId, logType mexadomain.CCLogType, logValue mexadomain.CCLogValue) error {
		test.AssertEqual(t, "casualty id should match", casualtyId, id)
		test.AssertEqual(t, "log type should be TreatStart", mexadomain.CCLogTypeTreatStart, logType)
		return nil
	}

	repos.Casualties.GetCasualtyByIdFunc = func(ctx context.Context, exerciseId int, id mexadomain.CasualtyId) (*mexadomain.Casualty, error) {
		return &mexadomain.Casualty{Id: id, ExerciseId: exerciseId, CaseId: 101}, nil
	}
	repos.CCLogs.GetLogsByCasualtyIdFunc = func(ctx context.Context, id mexadomain.CasualtyId) ([]mexadomain.CCLog, error) {
		return []mexadomain.CCLog{{Type: mexadomain.CCLogTypeTreatStart}}, nil
	}
	repos.Deterioration.GetDeteriorationByCasualtyFunc = func(ctx context.Context, id mexadomain.CasualtyId) ([]mexadomain.CadetDeterioration, error) {
		return nil, nil
	}

	u := chatdomain.Update{
		CallbackQuery: &struct {
			Id      string          `json:"id"`
			From    chatdomain.User `json:"from"`
			Message struct {
				chatdomain.Message

				Date        int                    `json:"date"`
				ReplyMarkup chatdomain.ReplyMarkup `json:"reply_markup"`
			} `json:"message"`
			ChatInstance string `json:"chat_instance"`
			Data         string `json:"data"`
		}{
			Data: treatStartPrefix + "::start:" + strconv.Itoa(casualtyId),
			From: chatdomain.User{Id: 123},
		},
	}
	u.CallbackQuery.Message.Chat.Id = 123
	u.CallbackQuery.Message.MessageId = 444

	err := svc.HandleCallbacks(ctx, u)
	test.NilErr(t, err)

	test.AssertEqual(t, "should have 1 reply", 1, len(bot.ReplyCalls))
	test.AssertEqual(t, "reply text should match", "Treatment started", bot.ReplyCalls[0].Text)
	test.AssertEqual(t, "should have 1 EditMessage call", 1, len(bot.EditMessageCalls))
}
