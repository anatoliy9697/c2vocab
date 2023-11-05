package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	SaveNewTgChat(*tcPkg.Chat) error
	StartState() (*tcPkg.State, error)
	TgChatByUserId(int) (*tcPkg.Chat, error)
	StateByCode(string) (*tcPkg.State, error)
	UpdateTgChat(*tcPkg.Chat) error
	CmdByCode(string) (*tcPkg.Cmd, error)
	AllExercises() []*tcPkg.Excersice
	ExcersiceByCode(string) (*tcPkg.Excersice, error)
}

func Init(c context.Context, p *pgxpool.Pool) (Repo, error) {
	return initPGRepo(c, p)
}
