package mexasqlite

import (
	"context"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/infra/db/sqlite"
	"mexa/internal/utils"
)

type UsersRepo struct {
	*sqlite.BaseRepo
}

func NewUsersRepo(base *sqlite.BaseRepo) (repo *UsersRepo) {
	return &UsersRepo{
		BaseRepo: base,
	}
}

func (repo *UsersRepo) CreateUserIfNotExists(ctx context.Context, id mexadomain.UserId, username string) (err error) {
	stmt := "INSERT INTO users (id, username, created_at) VALUES (?,?,?) ON CONFLICT DO NOTHING"
	_, err = repo.Exec(ctx, stmt, id, username, utils.TNow())
	return err
}
