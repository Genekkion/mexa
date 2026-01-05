package mexaservice

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
)

func (s *Service) genericErrorReply(ctx context.Context, u chatdomain.Update, err error) error {
	fmt.Println("Generic error:", err)
	return s.bot.Reply(ctx, u.ChatId(), "Something went wrong, please try again")
}
