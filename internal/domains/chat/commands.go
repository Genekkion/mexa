package chatdomain

import "context"

type Handler func(ctx context.Context, u Update) (err error)

type Command struct {
	Text        string  `json:"command"`
	Description string  `json:"description"`
	Handler     Handler `json:"-"`
}
