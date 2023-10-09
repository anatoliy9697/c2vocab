package repo

import (
	"context"
	"errors"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) ToInnerTgChat(u *usrPkg.User, outerChat *tgbotapi.Chat) (tc *tcPkg.TgChat, err error) {
	if tc, err = r.TgChatByUserId(u.Id); err == nil && tc == nil {
		state, _ := r.StartState()
		tc = tcPkg.NewTgChat(outerChat, u.Id, state)
		err = r.SaveNewTgChat(tc)
	}

	return tc, err
}

func (r pgRepo) SaveNewTgChat(tc *tcPkg.TgChat) error {
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
		tc.State.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) TgChatByUserId(usrId int32) (*tcPkg.TgChat, error) {
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var state *tcPkg.State
	state, err = r.StateByCode(stateCode)
	if err != nil {
		return nil, err
	}

	return &tcPkg.TgChat{TgId: tgId, UserId: usrId, State: state}, nil
}
