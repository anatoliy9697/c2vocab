package repo

import (
	"context"

	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	SaveNew(u *tgchat.TgChat) error
	StartState() (*tgchat.State, error)
	TgChatByUserId(int32) (*tgchat.TgChat, error)
	StateByCode(string) (*tgchat.State, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
