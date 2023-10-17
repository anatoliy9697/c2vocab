package repo

import (
	"context"

	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	SaveNewWL(*wlPkg.WordList) error
	ActiveWLByOwnerId(int32) ([]*wlPkg.WordList, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
