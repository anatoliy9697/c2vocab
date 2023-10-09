package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	ToInnerTgChat(*usrPkg.User, *tgbotapi.Chat) (*tcPkg.TgChat, error)
	SaveNewTgChat(*tcPkg.TgChat) error
	StartState() (*tcPkg.State, error)
	TgChatByUserId(int32) (*tcPkg.TgChat, error)
	StateByCode(string) (*tcPkg.State, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
