package repo

import (
	"context"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
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
		INSERT INTO c2v_word_list (active, name, frgn_lang_code, ntv_lang_code, owner_id, created_at)
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

func (r pgRepo) ActiveWLByOwnerId(ownerId int32) ([]*wlPkg.WordList, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT id, name, frgn_lang_code, ntv_lang_code, created_at
		FROM c2v_word_list
		WHERE owner_id = $1 AND active IS TRUE
	`
	rows, err := conn.Query(r.ctx, sql, ownerId)
	if err != nil {
		return nil, err
	}

	wls := make([]*wlPkg.WordList, 0)
	var id int32
	var name, frgnLangCode, ntvLangCode string
	var createdAt time.Time
	for rows.Next() {
		if err = rows.Scan(&id, &name, &frgnLangCode, &ntvLangCode, &createdAt); err != nil {
			return nil, err
		}
		wls = append(wls, &wlPkg.WordList{
			Id:        id,
			Active:    true,
			Name:      name,
			FrgnLang:  commons.LangByCode(frgnLangCode),
			NtvLang:   commons.LangByCode(ntvLangCode),
			OwnerId:   ownerId,
			CreatedAt: createdAt,
		})
	}

	return wls, nil
}

func (r pgRepo) WLById(id int32) (*wlPkg.WordList, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT active, name, frgn_lang_code, ntv_lang_code, owner_id, created_at
		FROM c2v_word_list
		WHERE id = $1
	`
	var active bool
	var name, frgnLangCode, ntvLangCode string
	var ownerId int32
	var createdAt time.Time
	err = conn.QueryRow(r.ctx, sql, id).Scan(
		&active,
		&name,
		&frgnLangCode,
		&ntvLangCode,
		&ownerId,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &wlPkg.WordList{
		Id:        id,
		Active:    active,
		Name:      name,
		FrgnLang:  commons.LangByCode(frgnLangCode),
		NtvLang:   commons.LangByCode(ntvLangCode),
		OwnerId:   ownerId,
		CreatedAt: createdAt,
	}, nil
}

func (r pgRepo) UpdateWL(wl *wlPkg.WordList) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		UPDATE c2v_word_list SET (active, name, frgn_lang_code, ntv_lang_code) = ($1, $2, $3, $4)
		WHERE id = $5
	`
	_, err = conn.Exec(r.ctx, sql,
		wl.Active,
		wl.Name,
		wl.FrgnLang.Code,
		wl.NtvLang.Code,
		wl.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) SaveNewWord(w *wlPkg.Word) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_word (frgn, ntv, wl_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var wId int32
	err = conn.QueryRow(r.ctx, sql,
		w.Foreign,
		w.Native,
		w.WLId,
		w.CreatedAt,
	).Scan(&wId)
	if err != nil {
		return err
	}

	w.Id = wId

	return nil
}
