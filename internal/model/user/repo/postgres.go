package repo

import (
	"context"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) IsExists(u *usrPkg.User) (bool, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	sql := `
		SELECT COUNT(*)
		FROM c2v_user
		WHERE tg_id=$1
	`
	var count int
	err = conn.QueryRow(r.ctx, sql, u.TgId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r pgRepo) SaveNew(u *usrPkg.User) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_user(tg_id, tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id
	`
	var usrId int
	err = conn.QueryRow(r.ctx, sql,
		u.TgId,
		u.TgUserName,
		u.TgFirstName,
		u.TgLastName,
		u.Lang.Code,
		u.TgIsBot,
	).Scan(&usrId)
	if err != nil {
		return err
	}

	u.Id = usrId

	return nil
}

func (r pgRepo) Update(u *usrPkg.User) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		UPDATE c2v_user SET (tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot) = ($1, $2, $3, $4, $5)
		WHERE tg_id = $6
		RETURNING id
	`
	var usrId int
	err = conn.QueryRow(r.ctx, sql,
		u.TgUserName,
		u.TgFirstName,
		u.TgLastName,
		u.Lang.Code,
		u.TgIsBot,
		u.TgId,
	).Scan(&usrId)
	if err != nil {
		return err
	}

	u.Id = usrId

	return nil
}

func (r pgRepo) ById(userId int) (u *usrPkg.User, err error) {
	var conn *pgxpool.Conn
	if conn, err = r.pool.Acquire(r.ctx); err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT
			tg_id
			, tg_username
			, tg_is_bot
			, tg_first_name
			, tg_last_name
			, tg_lang_code
		FROM c2v_user
		WHERE id = $1
	`
	var tgId int
	var tgUserName, tgFirstName, tgLastName, tgLangCode string
	var tgIsBot bool
	err = conn.QueryRow(r.ctx, sql, userId).Scan(
		&tgId,
		&tgUserName,
		&tgIsBot,
		&tgFirstName,
		&tgLastName,
		&tgLangCode,
	)
	if err != nil {
		return nil, err
	}

	u = &usrPkg.User{
		Id:          userId,
		TgId:        tgId,
		TgUserName:  tgUserName,
		TgFirstName: tgFirstName,
		TgLastName:  tgLastName,
		Lang:        commons.LangByCode(tgLangCode),
		TgIsBot:     tgIsBot,
	}

	return u, nil
}
