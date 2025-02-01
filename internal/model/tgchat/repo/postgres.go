package repo

import (
	"context"
	"errors"
	"strconv"
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

	if err := initExercisesTaskTextTmpls(); err != nil {
		return nil, err
	}

	return &pgRepo{c, p}, nil
}

func (r pgRepo) IsChatExistsByUserId(userId int) (bool, error) {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	sql := `
		SELECT COUNT(*)
		FROM c2v_tg_chat
		WHERE user_id = $1
	`
	var count int
	err = conn.QueryRow(r.c, sql, userId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r pgRepo) SaveNewChat(tc *tcPkg.Chat, handlerCode string) error {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_tg_chat(tg_id, user_id, state_code, usr_last_act_dt, created_at, worker, in_work_from)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4, CURRENT_TIMESTAMP)
	`
	_, err = conn.Exec(r.c, sql,
		tc.TgId,
		tc.UserId,
		tc.State.Code,
		handlerCode,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) ChatByUserId(usrId int) (*tcPkg.Chat, error) {
	conn, err := r.p.Acquire(r.c)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var tgId, wlId, wordId, botLastMsgId int
	var stateCode, wlFrgnLangCode, wlNtvLangCode, wordFrgn, excersiceCode, trainedWordsIds string
	sql := `
		SELECT
			tg_id
			, state_code
			, COALESCE(wl_frgn_lang_code, '')
			, COALESCE(wl_ntv_lang_code, '')
			, COALESCE(wl_id, 0)
			, COALESCE(word_frgn, '')
			, COALESCE(word_id, 0)
			, COALESCE(exercise_code, '')
			, COALESCE(trained_words_ids, '')
			, COALESCE(bot_last_msg_id, 0)
		FROM c2v_tg_chat
		WHERE user_id = $1
	`
	err = conn.QueryRow(r.c, sql, usrId).Scan(
		&tgId,
		&stateCode,
		&wlFrgnLangCode,
		&wlNtvLangCode,
		&wlId,
		&wordFrgn,
		&wordId,
		&excersiceCode,
		&trainedWordsIds,
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
		TgId:            tgId,
		UserId:          usrId,
		State:           state,
		WLFrgnLang:      commons.LangByCode(wlFrgnLangCode),
		WLNtvLang:       commons.LangByCode(wlNtvLangCode),
		WLId:            wlId,
		WordFrgn:        wordFrgn,
		WordId:          wordId,
		ExcersiceCode:   excersiceCode,
		TrainedWordsIds: trainedWordsIds,
		BotLastMsgId:    botLastMsgId,
	}, nil
}

func (r pgRepo) UpdateChat(tc *tcPkg.Chat, usrActivity bool) error {
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
		UPDATE c2v_tg_chat SET (
			tg_id
			, state_code
			, usr_last_act_dt
			, wl_frgn_lang_code
			, wl_ntv_lang_code
			, wl_id
			, word_frgn
			, word_id
			, words_ids
			, exercise_code
			, trained_words_ids
			, bot_last_msg_id
		) = (
			$1
			, $2
			, CASE WHEN $3 IS TRUE THEN CURRENT_TIMESTAMP ELSE usr_last_act_dt END
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
		)
		WHERE user_id = $12
	`
	_, err = conn.Exec(r.c, sql,
		tc.TgId,
		tc.State.Code,
		usrActivity,
		wlFrgnLangCode,
		wlNtvLangCode,
		tc.WLId,
		tc.WordFrgn,
		tc.WordId,
		tc.WordsIdsStr(),
		tc.ExcersiceCode,
		tc.TrainedWordsIds,
		tc.BotLastMsgId,
		tc.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) UnlockChatByUserId(userId int) (err error) {
	var conn *pgxpool.Conn
	if conn, err = r.p.Acquire(r.c); err != nil {
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
	_, err = conn.Exec(r.c, sql, userId)

	return err
}

func (r pgRepo) LockChatByUserId(userId int, handlerCode string, timeForReassign int, lockAttemptsAmount int, timeForNextLockAttempt int) (err error) {
	var conn *pgxpool.Conn
	if conn, err = r.p.Acquire(r.c); err != nil {
		return err
	}
	defer conn.Release()

	var locked bool
	sql := `
		UPDATE c2v_tg_chat
		SET
			worker = CASE
				WHEN worker IS NULL OR in_work_from + interval '` + strconv.Itoa(timeForReassign) + ` seconds' <= CURRENT_TIMESTAMP
				THEN $1
				ELSE worker
			END
			, in_work_from = CASE
				WHEN worker IS NULL OR in_work_from + interval '` + strconv.Itoa(timeForReassign) + ` seconds' <= CURRENT_TIMESTAMP
				THEN CURRENT_TIMESTAMP
				ELSE in_work_from
			END
		WHERE
			user_id = $2
		RETURNING
			CASE
				WHEN worker = $1 THEN true ELSE false 
			END
	`
	for lockAttemptsCount := 1; lockAttemptsCount <= lockAttemptsAmount; lockAttemptsCount++ {
		if err = conn.QueryRow(r.c, sql, handlerCode, userId).Scan(&locked); err != nil {
			return err
		}
		if locked {
			return nil
		}
		time.Sleep(time.Duration(timeForNextLockAttempt) * time.Millisecond)
	}

	return errors.New("lock attempts limit reached")
}
