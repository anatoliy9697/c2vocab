package repo

import (
	"context"

	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) ToInner(outerU *tgbotapi.User) (u *usrPkg.User, err error) {
	u = usrPkg.MapToInner(outerU)

	var userExists bool
	if userExists, err = r.IsExists(u); err == nil {
		if userExists {
			err = r.Update(u)
		} else {
			err = r.SaveNew(u)
		}
	}

	return u, err
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
		INSERT INTO c2v_user(tg_id, tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var usrId int32
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
	var usrId int32
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
