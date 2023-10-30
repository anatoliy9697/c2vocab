package repo

import (
	"context"

	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	SaveNewWL(*wlPkg.WordList) error
	ActiveWLByOwnerId(int) ([]*wlPkg.WordList, error)
	WLById(int) (*wlPkg.WordList, error)
	UpdateWL(*wlPkg.WordList) error
	SaveNewWord(*wlPkg.Word) error
	ActiveWordsByWLId(int) ([]*wlPkg.Word, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
