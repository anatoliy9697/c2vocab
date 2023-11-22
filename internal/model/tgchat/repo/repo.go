package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	IsChatExistsByUserId(int) (bool, error)
	SaveNewChat(*tcPkg.Chat, string) error
	StartState() (*tcPkg.State, error)
	ChatByUserId(int) (*tcPkg.Chat, error)
	StateByCode(string) (*tcPkg.State, error)
	UpdateChat(*tcPkg.Chat, bool) error
	CmdByCode(string) (*tcPkg.Cmd, error)
	AllExercises() []*tcPkg.Excersice
	ExcersiceByCode(string) (*tcPkg.Excersice, error)
	UnlockChatByUserId(int) error
	LockChatByUserId(int, string, int, int, int) error
}

func Init(c context.Context, p *pgxpool.Pool) (Repo, error) {
	return initPGRepo(c, p)
}
