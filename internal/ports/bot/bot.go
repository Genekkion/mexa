package botports

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
)

type Bot interface {
	Reply(ctx context.Context, chatId chatdomain.ChatId, text string, opts ...chatdomain.SendOpt) (err error)
	EditMessage(ctx context.Context, chatId chatdomain.ChatId, messageId chatdomain.MessageId, text string,
		opts ...chatdomain.SendOpt) (err error)
	SetupCommands(ctx context.Context, commands []chatdomain.Command) (err error)

	GetUpdates() (u []chatdomain.Update, err error)
	UpdateOffset(offset int)
}
