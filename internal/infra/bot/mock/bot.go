package botmock

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
)

type Bot struct {
	ReplyCalls         []ReplyCall
	EditMessageCalls   []EditMessageCall
	SetupCommandsCalls [][]chatdomain.Command
}

type ReplyCall struct {
	ChatId chatdomain.ChatId
	Text   string
	Opts   []chatdomain.SendOpt
}

type EditMessageCall struct {
	ChatId    chatdomain.ChatId
	MessageId chatdomain.MessageId
	Text      string
	Opts      []chatdomain.SendOpt
}

func (b *Bot) Reply(ctx context.Context, chatId chatdomain.ChatId, text string, opts ...chatdomain.SendOpt) error {
	b.ReplyCalls = append(b.ReplyCalls, ReplyCall{ChatId: chatId, Text: text, Opts: opts})
	return nil
}

func (b *Bot) EditMessage(ctx context.Context, chatId chatdomain.ChatId, messageId chatdomain.MessageId, text string, opts ...chatdomain.SendOpt) error {
	b.EditMessageCalls = append(b.EditMessageCalls, EditMessageCall{ChatId: chatId, MessageId: messageId, Text: text, Opts: opts})
	return nil
}

func (b *Bot) SetupCommands(ctx context.Context, commands []chatdomain.Command) error {
	b.SetupCommandsCalls = append(b.SetupCommandsCalls, commands)
	return nil
}

func (b *Bot) GetUpdates() ([]chatdomain.Update, error) {
	return nil, nil
}

func (b *Bot) UpdateOffset(offset int) {}
