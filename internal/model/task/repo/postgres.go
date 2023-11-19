package repo

import (
	"context"

	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) Tasks(batchSize int) ([]tskPkg.Task, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT
			user_id
			, 'to_main_menu' AS  task_type
		FROM c2v_tg_chat
		WHERE state_code <> 'main_menu' AND usr_last_act_dt + interval '1 minute' <= CURRENT_TIMESTAMP
		LIMIT $1
	`
	rows, err := conn.Query(r.ctx, sql, batchSize)
	if err != nil {
		return nil, err
	}

	tasks := make([]tskPkg.Task, 0, batchSize)
	var userId int
	var taskType string
	for rows.Next() {
		if err = rows.Scan(&userId, &taskType); err != nil {
			return nil, err
		}
		tasks = append(tasks, tskPkg.Task{
			Type:   taskType,
			UserId: userId,
		})
	}

	return tasks, nil
}
