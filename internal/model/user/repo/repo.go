package repo

import (
	"context"

	"github.com/anatoliy9697/c2vocab/internal/model/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Set(u *user.User) (bool, error)
	IsExists(u *user.User) (bool, error)
	SaveNew(u *user.User) error
	Update(u *user.User) error
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
