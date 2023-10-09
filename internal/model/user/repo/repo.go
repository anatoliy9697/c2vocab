package repo

import (
	"context"

	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	ToInner(*tgbotapi.User) (*usrPkg.User, error)
	IsExists(*usrPkg.User) (bool, error)
	SaveNew(*usrPkg.User) error
	Update(*usrPkg.User) error
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
