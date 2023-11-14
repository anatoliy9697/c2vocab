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

func (r pgRepo) WLByIdAndUserId(wlId, userId int) (*wlPkg.WordList, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	mem_percentage_formula := `
		power(
			2.718
			, -1 * (EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - last_cor_ans_dt)) / 60) / (2078 * power(3, fc_trainings_num - 1))
		) * 100
	`
	mem_percentage_sql := `
		SELECT 
			ROUND(AVG(CASE
				WHEN ws.fc_trainings_num IS NOT NULL AND ws.last_cor_ans_dt IS NOT NULL
				THEN ` + mem_percentage_formula + `
				ELSE 0
			END))::integer
		FROM c2v_word w
		LEFT JOIN c2v_word_stat ws ON ws.word_id = w.id AND ws.user_id = $2
		WHERE w.wl_id = $1 AND w.active IS TRUE
	`
	sql := `
		SELECT
			active
			, name
			, frgn_lang_code
			, ntv_lang_code
			, (SELECT COUNT(*) FROM c2v_word WHERE wl_id = $1 AND active IS TRUE) AS words_num
			, COALESCE((` + mem_percentage_sql + `), 0) AS mem_percentage
			, owner_id
			, created_at
		FROM c2v_word_list
		WHERE id = $1
	`
	var active bool
	var name, frgnLangCode, ntvLangCode string
	var wordsNum, ownerId, memPercentage int
	var createdAt time.Time
	err = conn.QueryRow(r.ctx, sql, wlId, userId).Scan(
		&active,
		&name,
		&frgnLangCode,
		&ntvLangCode,
		&wordsNum,
		&memPercentage,
		&ownerId,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &wlPkg.WordList{
		Id:            wlId,
		Active:        active,
		Name:          name,
		FrgnLang:      commons.LangByCode(frgnLangCode),
		NtvLang:       commons.LangByCode(ntvLangCode),
		WordsNum:      wordsNum,
		MemPercentage: memPercentage,
		OwnerId:       ownerId,
		CreatedAt:     createdAt,
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

func (r pgRepo) WordByIdAndUserId(wordId, userId int) (*wlPkg.Word, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	word_stat_formula := `
		ROUND(power(
			2.718
			, -1 * (EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - last_cor_ans_dt)) / 60) / (2078 * power(3, fc_trainings_num - 1))
		) * 100)::integer
	`
	word_stat_sql := `
		SELECT
			CASE
				WHEN fc_trainings_num IS NOT NULL AND last_cor_ans_dt IS NOT NULL THEN ` + word_stat_formula + ` ELSE 0
			END
		FROM c2v_word_stat
		WHERE word_id = $1 AND user_id = $2
	`
	sql := `
		SELECT
			active
			, frgn
			, ntv
			, COALESCE((` + word_stat_sql + `), 0) AS mem_percentage
			, created_at
		FROM c2v_word
		WHERE id = $1
	`
	var active bool
	var frgn, ntv string
	var createdAt time.Time
	var memPercentage int
	err = conn.QueryRow(r.ctx, sql, wordId, userId).Scan(
		&active,
		&frgn,
		&ntv,
		&memPercentage,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &wlPkg.Word{
		Id:            wordId,
		Active:        active,
		Foreign:       frgn,
		Native:        ntv,
		MemPercentage: memPercentage,
		CreatedAt:     createdAt,
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

func (r pgRepo) IsWordStatExists(wordId, userId int) (bool, error) {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	sql := `
		SELECT COUNT(*)
		FROM c2v_word_stat
		WHERE word_id = $1 AND user_id = $2
	`
	var count int
	err = conn.QueryRow(r.ctx, sql, wordId, userId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r pgRepo) CreateWordStat(wordId, userId int, isAnswerCorrect bool) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO c2v_word_stat (
			word_id
			, last_training_dt
			, trainings_num
			, correct_answers_num
			, last_cor_ans_dt
			, fc_trainings_num
			, fc_next_training_dt
			, user_id
		) VALUES (
			$1
			, CURRENT_TIMESTAMP
			, 1
			, CASE WHEN $2 IS TRUE THEN 1 ELSE 0 END
			, CASE WHEN $2 IS TRUE THEN CURRENT_TIMESTAMP ELSE NULL END
			, CASE WHEN $2 IS TRUE THEN 1 ELSE NULL END
			, CASE WHEN $2 IS TRUE THEN CURRENT_TIMESTAMP + interval '180 minutes' ELSE NULL END
			, $3
		)
	`
	_, err = conn.Exec(r.ctx, sql, wordId, isAnswerCorrect, userId)

	return err
}

func (r pgRepo) RegistrateWordTraining(wordId, userId int, isAnswerCorrect bool, memPercentageDowngrade int) error {
	conn, err := r.pool.Acquire(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	fc_trainings_num := `
		CASE
			WHEN $2 IS TRUE AND (fc_next_training_dt IS NULL OR $4 = -1 OR COALESCE(fc_trainings_num, 0) - $4 <= 0)
				THEN 1
			WHEN $2 IS TRUE AND CURRENT_TIMESTAMP >= fc_next_training_dt
				THEN COALESCE(fc_trainings_num, 0) - $4 + 1
			WHEN $2 IS TRUE
				THEN COALESCE(fc_trainings_num, 0) - $4
			ELSE
				fc_trainings_num
		END
	`
	fc_next_training_dt := `
		CASE
			WHEN $2 IS TRUE AND (fc_next_training_dt IS NULL OR $4 = -1 OR COALESCE(fc_trainings_num, 0) - $4 <= 0)
				THEN CURRENT_TIMESTAMP + (interval '1 minute' * 180)
			WHEN $2 IS TRUE AND CURRENT_TIMESTAMP >= fc_next_training_dt
				THEN CURRENT_TIMESTAMP + (interval '1 minute' * 180 * power(3, COALESCE(fc_trainings_num, 0) - $4))
			WHEN $2 IS TRUE
				THEN CURRENT_TIMESTAMP + (interval '1 minute' * 180 * power(3, COALESCE(fc_trainings_num, 0) - $4 - 1))
			ELSE
				fc_next_training_dt
		END
	`
	sql := `
		UPDATE c2v_word_stat
		SET (
			last_training_dt
			, trainings_num
			, correct_answers_num
			, last_cor_ans_dt
			, fc_trainings_num
			, fc_next_training_dt
		) = (
			CURRENT_TIMESTAMP
			, trainings_num + 1
			, CASE WHEN $2 IS TRUE THEN correct_answers_num + 1 ELSE correct_answers_num END
			, CASE WHEN $2 IS TRUE THEN CURRENT_TIMESTAMP ELSE last_cor_ans_dt END
			, ` + fc_trainings_num + `
			, ` + fc_next_training_dt + `
		) WHERE word_id = $1 AND user_id = $3
	`
	_, err = conn.Exec(r.ctx, sql, wordId, isAnswerCorrect, userId, memPercentageDowngrade)

	return err
}
