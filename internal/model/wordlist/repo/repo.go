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
	WordById(int) (*wlPkg.Word, error)
	UpdateWord(*wlPkg.Word) error
	NextWordForTraining(int, string) (*wlPkg.Word, error)
	WordSelectionAnswerOptions(*wlPkg.Word, bool, string, int, int) ([]wlPkg.AnswerOption, error)
	IsWordStatExists(int, int) (bool, error)
	CreateWordStat(int, int, bool) error
	RegistrateWordTraining(int, int, bool) error
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
