package mexaservice

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	botmock "mexa/internal/infra/bot/mock"
	"mexa/internal/infra/db/sqlite"
	mexasqlite "mexa/internal/infra/db/sqlite/mexa"
	fsmmock "mexa/internal/infra/fsm/mock"
	fsmports "mexa/internal/ports/fsm"
	"mexa/internal/test"
	"strconv"
	"testing"
)

func TestService_Integration_Scenario(t *testing.T) {
	db := sqlite.NewTestDb(t)
	defer db.Close()
	baseRepo := sqlite.NewBaseRepo(db)

	repos := Repos{
		Transactional: &baseRepo.Transactional,
		Users:         mexasqlite.NewUsersRepo(&baseRepo),
		Cases:         mexasqlite.NewCasesRepo(&baseRepo),
		Casualties:    mexasqlite.NewCasualtiesRepo(&baseRepo),
		Exercises:     mexasqlite.NewExercisesRepo(&baseRepo),
		Deterioration: mexasqlite.NewCasualtiesDeteriorationRepo(&baseRepo),
		ExLogs:        mexasqlite.NewExLogsRepo(&baseRepo),
		CCLogs:        mexasqlite.NewCadetCaseLogsRepo(&baseRepo),
	}

	bot := &botmock.Bot{}
	fsm := fsmmock.New()
	adminId := chatdomain.UserId(123)

	ctx := context.Background()

	// 1. Add Exercise and Case to DB
	exId, err := repos.Exercises.AddExercise(ctx, "EX1", "Exercise 1")
	test.NilErr(t, err)

	_, err = repos.Cases.AddCase(ctx, *exId, mexadomain.CaseValue{
		Summary: "Test Case",
	})
	test.NilErr(t, err)

	config := ServiceConfig{
		Bot:      bot,
		Admins:   []chatdomain.UserId{adminId},
		Fsm:      fsm,
		Repos:    repos,
		Exercise: mexadomain.Exercise{Id: *exId, Code: "EX1", Name: "Exercise 1"},
	}

	svc, err := NewService(ctx, config)
	test.NilErr(t, err)

	// --- SCENARIO START ---

	// 2. Admin starts exercise
	u := chatdomain.Update{
		Message: &chatdomain.Message{
			From: chatdomain.User{Id: int(adminId), Username: "admin"},
			Chat: struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			}{Id: int(adminId)},
			Text: "/ex_start",
		},
	}
	err = svc.HandleCommands(ctx, u)
	test.NilErr(t, err)
	test.AssertEqual(t, "Exercise should be started", fsmports.StateExStarted, fsm.FsmState())

	// 3. Admin attaches a case to a cadet
	// Step 3a: /attach command
	u.Message.Text = "/attach"
	err = svc.HandleCommands(ctx, u)
	test.NilErr(t, err)
	test.AssertEqual(t, "User should be in AttachingCase state", fsmports.UserStateAttachingCase, fsm.UserState(adminId))

	// Step 3b: Enter 4D number
	u.Message.Text = "1234"
	err = svc.HandleText(ctx, u)
	test.NilErr(t, err)
	test.Assert(t, "Should have received cases list", len(bot.ReplyCalls) > 0)
	lastReply := bot.ReplyCalls[len(bot.ReplyCalls)-1]
	test.Assert(t, "Reply should contain '1234'", contains(lastReply.Text, "1234"))

	// Step 3c: Select case from inline keyboard
	caseId := 1
	u = chatdomain.Update{
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
			Data: "attach::select:" + strconv.Itoa(caseId),
			From: chatdomain.User{Id: int(adminId)},
		},
	}
	u.CallbackQuery.Message.Chat.Id = int(adminId)
	u.CallbackQuery.Message.Text = lastReply.Text // Important for regex to find 4D

	err = svc.HandleCallbacks(ctx, u)
	test.NilErr(t, err)
	test.AssertEqual(t, "User should be back to Default state", fsmports.UserStateDefault, fsm.UserState(adminId))

	// Verify casualty added to DB
	cas, err := repos.Casualties.GetCasualtyBy4D(ctx, *exId, 1234)
	test.NilErr(t, err)
	test.AssertEqual(t, "CaseId should match", caseId, cas.CaseId)

	// 4. Start Treatment
	u.CallbackQuery.Data = "treat_start::start:" + strconv.Itoa(cas.Id)
	err = svc.HandleCallbacks(ctx, u)
	test.NilErr(t, err)

	// Verify log in DB
	logs, err := repos.CCLogs.GetLogsByCasualtyId(ctx, cas.Id)
	test.NilErr(t, err)
	test.AssertEqual(t, "Should have 1 log", 1, len(logs))
	test.AssertEqual(t, "Log should be TreatStart", mexadomain.CCLogTypeTreatStart, logs[0].Type)

	// 5. End Treatment (Success)
	u.CallbackQuery.Data = "treat_end::result:" + strconv.Itoa(cas.Id) + ":success"
	err = svc.HandleCallbacks(ctx, u)
	test.NilErr(t, err)

	// Verify second log in DB
	logs, err = repos.CCLogs.GetLogsByCasualtyId(ctx, cas.Id)
	test.NilErr(t, err)
	test.AssertEqual(t, "Should have 2 logs", 2, len(logs))
	test.AssertEqual(t, "Second log should be TreatEnd", mexadomain.CCLogTypeTreatEnd, logs[1].Type)
	test.AssertEqual(t, "Outcome should be Success", mexadomain.CCLogEndOutcomeSuccess, *logs[1].Value.Outcome)
}

func contains(s, substr string) bool {
	return (len(s) >= len(substr)) && (s == substr || (len(substr) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
