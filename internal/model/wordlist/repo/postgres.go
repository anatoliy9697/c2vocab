package repo

import (
	"context"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

var reserveAnswerOptions = map[string][]wlPkg.AnswerOption{
	"en": {{Answer: "reserve", IsCorrect: "0"}, {Answer: "explain", IsCorrect: "0"}, {Answer: "example", IsCorrect: "0"}},
	"ru": {{Answer: "запасной", IsCorrect: "0"}, {Answer: "впечатлять", IsCorrect: "0"}, {Answer: "произношение", IsCorrect: "0"}},
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
	var wlId int
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

func (r pgRepo) ActiveWLByOwnerId(ownerId int) ([]*wlPkg.WordList, error) {
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
	var id int
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

func (r pgRepo) WLById(id int) (*wlPkg.WordList, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT
			active
			, name
			, frgn_lang_code
			, ntv_lang_code
			, (SELECT COUNT(*) FROM c2v_word WHERE wl_id = $1 AND active IS TRUE) AS words_num
			, owner_id
			, created_at
		FROM c2v_word_list
		WHERE id = $1
	`
	var active bool
	var name, frgnLangCode, ntvLangCode string
	var wordsNum, ownerId int
	var createdAt time.Time
	err = conn.QueryRow(r.ctx, sql, id).Scan(
		&active,
		&name,
		&frgnLangCode,
		&ntvLangCode,
		&wordsNum,
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
		WordsNum:  wordsNum,
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
		INSERT INTO c2v_word (frgn, ntv, wl_id, active, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var wId int
	err = conn.QueryRow(r.ctx, sql,
		w.Foreign,
		w.Native,
		w.WLId,
		true,
		w.CreatedAt,
	).Scan(&wId)
	if err != nil {
		return err
	}

	w.Id = wId

	return nil
}

func (r pgRepo) ActiveWordsByWLId(wlId int) ([]*wlPkg.Word, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT id, frgn, ntv, created_at
		FROM c2v_word
		WHERE wl_id = $1 AND active IS TRUE
	`
	rows, err := conn.Query(r.ctx, sql, wlId)
	if err != nil {
		return nil, err
	}

	words := make([]*wlPkg.Word, 0)
	var id int
	var frgn, ntv string
	var createdAt time.Time
	for rows.Next() {
		if err = rows.Scan(&id, &frgn, &ntv, &createdAt); err != nil {
			return nil, err
		}
		words = append(words, &wlPkg.Word{
			Id:        id,
			Foreign:   frgn,
			Native:    ntv,
			WLId:      wlId,
			Active:    true,
			CreatedAt: createdAt,
		})
	}

	return words, nil
}

func (r pgRepo) WordById(id int) (*wlPkg.Word, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `
		SELECT
			active
			, frgn
			, ntv
			, created_at
		FROM c2v_word
		WHERE id = $1
	`
	var active bool
	var frgn, ntv string
	var createdAt time.Time
	err = conn.QueryRow(r.ctx, sql, id).Scan(
		&active,
		&frgn,
		&ntv,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &wlPkg.Word{
		Id:        id,
		Active:    active,
		Foreign:   frgn,
		Native:    ntv,
		CreatedAt: createdAt,
	}, nil
}

func (r pgRepo) UpdateWord(w *wlPkg.Word) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		UPDATE c2v_word SET (active, frgn, ntv) = ($1, $2, $3)
		WHERE id = $4
	`
	_, err = conn.Exec(r.ctx, sql,
		w.Active,
		w.Foreign,
		w.Native,
		w.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r pgRepo) NextWordForTraining(wlId int, excludedIds string) (*wlPkg.Word, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	optionalCondition := ""
	if excludedIds != "" {
		optionalCondition = " AND id NOT IN (" + excludedIds + ")"
	}

	sql := `
		SELECT
			id
			, frgn
			, ntv
			, created_at
		FROM c2v_word
		WHERE wl_id = $1 AND active IS TRUE` + optionalCondition + `
		ORDER BY RANDOM()
		LIMIT 1
	`
	var id int
	var frgn, ntv string
	var createdAt time.Time
	err = conn.QueryRow(r.ctx, sql, wlId).Scan(
		&id,
		&frgn,
		&ntv,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &wlPkg.Word{
		Id:        id,
		Active:    true,
		Foreign:   frgn,
		Native:    ntv,
		CreatedAt: createdAt,
	}, nil
}

func (r pgRepo) WordSelectionAnswerOptions(word *wlPkg.Word, frgnOpts bool, langCode string, userId, optsLimit int) (opts []wlPkg.AnswerOption, err error) {
	opts = make([]wlPkg.AnswerOption, 0, optsLimit)

	defer func() {
		w := ""
		if frgnOpts {
			w = word.Foreign
		} else {
			w = word.Native
		}
		opts = append(opts, wlPkg.AnswerOption{Answer: w, IsCorrect: "1"})
		if len(opts) < optsLimit {
			opts = append(opts, reserveAnswerOptions[langCode][0:optsLimit-len(opts)]...)
		}
		opts = wlPkg.MixAnswerOptions(opts)
	}()

	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	sql := `
		SELECT
			CASE
				WHEN wl.frgn_lang_code = $2 THEN w.frgn
				WHEN wl.ntv_lang_code = $2 THEN w.ntv
			END AS word
		FROM c2v_word w
		JOIN c2v_word_list wl ON w.wl_id = wl.id
		WHERE
			wl.owner_id = $1
			AND w.active IS TRUE
			AND wl.active IS TRUE
			AND w.id <> $3
			AND (wl.frgn_lang_code = $2 OR wl.ntv_lang_code = $2)
		ORDER BY RANDOM()
		LIMIT $4
	` // TODO: В будущем надо оптимизировать
	var rows pgx.Rows
	var otherW string
	if rows, err = conn.Query(r.ctx, sql, userId, langCode, word.Id, optsLimit-1); err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&otherW); err != nil {
			continue
		}
		opts = append(opts, wlPkg.AnswerOption{Answer: otherW, IsCorrect: "0"})
	}

	return
}
