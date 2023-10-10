package repo

import (
	"context"

	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	IsExists(*usrPkg.User) (bool, error)
	SaveNew(*usrPkg.User) error
	Update(*usrPkg.User) error
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
