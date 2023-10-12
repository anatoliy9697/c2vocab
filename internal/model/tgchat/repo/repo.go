package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	SaveNewTgChat(*tcPkg.TgChat) error
	StartState() (*tcPkg.State, error)
	TgChatByUserId(int32) (*tcPkg.TgChat, error)
	StateByCode(string) (*tcPkg.State, error)
	UpdateTgChat(*tcPkg.TgChat) error
	CmdByCode(string) (*tcPkg.Cmd, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
