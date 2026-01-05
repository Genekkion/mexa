package mexaports

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
)

// UsersRepo is a repository for users.
type UsersRepo interface {
	CreateUserIfNotExists(ctx context.Context, id mexadomain.UserId, username string) (err error)
}
