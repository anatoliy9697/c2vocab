package repo

import (
	"context"
	"errors"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) (*pgRepo, error) {
	if err := initStateMsgTmpls(); err != nil {
		return nil, err
	}

	return &pgRepo{c, p}, nil
}

func (r pgRepo) SaveNewTgChat(tc *tcPkg.Chat) error {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_tg_chat(tg_id, user_id, state_code, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = conn.Exec(r.c, sql,
		tc.TgId,
		tc.UserId,
		tc.State.Code,
		tc.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) TgChatByUserId(usrId int) (*tcPkg.Chat, error) {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var tgId, wlId, botLastMsgId int
	var createdAt time.Time
	var stateCode, wlFrgnLangCode, wlNtvLangCode, wordFrgn string
	sql := `
		SELECT
			tg_id
			, created_at
			, state_code
			, COALESCE(wl_frgn_lang_code, '')
			, COALESCE(wl_ntv_lang_code, '')
			, COALESCE(wl_id, 0)
			, COALESCE(word_frgn, '')
			, COALESCE(bot_last_msg_id, 0)
		FROM c2v_tg_chat
		WHERE user_id = $1
	`
	err = conn.QueryRow(r.c, sql, usrId).Scan(
		&tgId,
		&createdAt,
		&stateCode,
		&wlFrgnLangCode,
		&wlNtvLangCode,
		&wlId,
		&wordFrgn,
		&botLastMsgId,
	)
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

	return &tcPkg.Chat{
		TgId:         tgId,
		CreatedAt:    createdAt,
		UserId:       usrId,
		State:        state,
		WLFrgnLang:   commons.LangByCode(wlFrgnLangCode),
		WLNtvLang:    commons.LangByCode(wlNtvLangCode),
		WLId:         wlId,
		WordFrgn:     wordFrgn,
		BotLastMsgId: botLastMsgId,
	}, nil
}

func (r pgRepo) UpdateTgChat(tc *tcPkg.Chat) error {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return err
	}
	defer conn.Release()

	var wlFrgnLangCode, wlNtvLangCode string
	if tc.WLFrgnLang != nil {
		wlFrgnLangCode = tc.WLFrgnLang.Code
	}
	if tc.WLNtvLang != nil {
		wlNtvLangCode = tc.WLNtvLang.Code
	}
	sql := `
		UPDATE c2v_tg_chat SET (tg_id, state_code, wl_frgn_lang_code, wl_ntv_lang_code, wl_id, word_frgn, bot_last_msg_id) = ($1, $2, $3, $4, $5, $6, $7)
		WHERE user_id = $8
	`
	_, err = conn.Exec(r.c, sql,
		tc.TgId,
		tc.State.Code,
		wlFrgnLangCode,
		wlNtvLangCode,
		tc.WLId,
		tc.WordFrgn,
		tc.BotLastMsgId,
		tc.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}
