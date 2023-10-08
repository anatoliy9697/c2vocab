package repo

import (
	"context"

	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) SaveNew(tc *tgchat.TgChat) error {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_tg_chat(tg_id, user_id, state_code)
		VALUES ($1, $2, $3)
	`
	_, err = conn.Exec(r.c, sql,
		tc.TgId,
		tc.UserId,
		tc.StateCode,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) TgChatByUserId(usrId int32) (*tgchat.TgChat, error) {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var tgId int64
	var stateCode string
	sql := `
		SELECT tg_id, state_code FROM c2v_tg_chat
		WHERE user_id = $1
	`
	err = conn.QueryRow(r.c, sql, usrId).Scan(&tgId, &stateCode)
	if err != nil {
		return nil, err
	}

	return &tgchat.TgChat{TgId: tgId, UserId: usrId, StateCode: stateCode}, nil
}
