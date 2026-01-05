package mexaservice

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
)

func (s *Service) initCommands(ctx context.Context) (err error) {
	commands := []chatdomain.Command{
		{
			Text:        "user",
			Description: "Get user info",
			Handler:     s.cmdUser,
		},
		{
			Text:        "batch",
			Description: "Get batch info",
			Handler:     s.cmdBatch,
		},
		{
			Text:        "exercise",
			Description: "Get exercise info",
			Handler:     s.cmdExercise,
		},

		{
			Text:        "cases",
			Description: "List available cases",
			Handler:     s.cmdListCases,
		},

		{
			Text:        "ex_start",
			Description: "Start the exercise",
			Handler:     s.wrapAdmin(s.cmdExStart),
		},
		{
			Text:        "ex_end",
			Description: "End the exercise",
			Handler:     s.wrapAdmin(s.cmdExEnd),
		},

		{
			Text:        "casualties",
			Description: "List all casualties",
			Handler:     s.cmdCasualties,
		},
		{
			Text:        "casualties_check",
			Description: "Check for casualty by 4D number",
			Handler:     s.cmdCasualtyCheck,
		},

		{
			Text:        "attach",
			Description: "Attach case to cadet to create a new casualty",
			Handler:     s.wrapExStarted(s.cmdAttach),
		},

		{
			Text:        "quit",
			Description: "Quit all current actions",
			Handler:     s.cmdQuit,
		},
	}

	s.commands = make(map[string]chatdomain.Command, len(commands))
	for _, cmd := range commands {
		s.commands[cmd.Text] = cmd
	}

	return s.bot.SetupCommands(ctx, commands)
}
