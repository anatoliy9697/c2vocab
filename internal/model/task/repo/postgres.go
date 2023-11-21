package repo

import (
	"context"
	"strconv"

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

func (r pgRepo) TasksWithLocking(handlerCode string, batchSize int, maxTimeForReassign int) ([]tskPkg.Task, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		WITH batch AS (
			SELECT
				user_id AS batch_user_id
				, 'to_main_menu' AS task_type
			FROM c2v_tg_chat
			WHERE 
				(COALESCE(worker, '') = '' OR in_work_from + interval '` + strconv.Itoa(maxTimeForReassign) + ` seconds' <= CURRENT_TIMESTAMP)
				AND (state_code <> 'main_menu' AND usr_last_act_dt + interval '10 minutes' <= CURRENT_TIMESTAMP)
			LIMIT $2
		)
		UPDATE c2v_tg_chat
		SET
			worker = $1,
			in_work_from = CURRENT_TIMESTAMP
		WHERE
			user_id IN (SELECT batch_user_id FROM batch)
		RETURNING
			user_id
			, (SELECT task_type FROM batch WHERE batch_user_id = user_id) AS task_type
	`
	rows, err := conn.Query(r.ctx, sql, handlerCode, batchSize)
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

func (r pgRepo) UnlockTaskByUserId(userId int) (err error) {
	var conn *pgxpool.Conn
	if conn, err = r.pool.Acquire(r.ctx); err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		UPDATE c2v_tg_chat
		SET
			worker = NULL,
			in_work_from = NULL
		WHERE
			user_id = $1
	`
	_, err = conn.Exec(r.ctx, sql, userId)

	return err
}
