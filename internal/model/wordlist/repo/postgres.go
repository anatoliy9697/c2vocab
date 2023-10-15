package repo

import (
	"context"

	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) SaveNewWL(wl *wlPkg.WordList) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_word_list(active, name, frgn_lang_code, ntv_lang_code, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var wlId int32
	err = conn.QueryRow(r.ctx, sql,
		wl.Active,
		wl.Name,
		wl.FrgnLang.Code,
		wl.NtvLang.Code,
		wl.OwnerId,
		wl.CreatedAt,
	).Scan(&wlId)
	if err != nil {
		return err
	}

	wl.Id = wlId

	return nil
}
