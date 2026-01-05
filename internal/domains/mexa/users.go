package mexadomain

import (
	chatdomain "mexa/internal/domains/chat"
	"time"
)

type UserId = chatdomain.UserId

type User struct {
	Id        UserId
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
