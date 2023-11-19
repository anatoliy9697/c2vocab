package repo

import (
	"context"

	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Tasks(int) ([]tskPkg.Task, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
