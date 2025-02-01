CREATE TABLE c2v_user (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	tg_id BIGINT UNIQUE NOT NULL,
	tg_username VARCHAR(32) UNIQUE NOT NULL,
	tg_first_name VARCHAR(64),
	tg_last_name VARCHAR(64),
	tg_lang_code VARCHAR(3) NOT NULL,
	tg_is_bot BOOLEAN NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE c2v_word_list (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	active BOOLEAN NOT NULL,
	name VARCHAR(2000) NOT NULL,
	frgn_lang_code VARCHAR(3) NOT NULL,
	ntv_lang_code VARCHAR(3) NOT NULL,
	owner_id BIGINT NOT NULL REFERENCES c2v_user(id),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE c2v_word (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	active BOOLEAN NOT NULL,
	frgn VARCHAR(2000) NOT NULL,
	ntv VARCHAR(2000) NOT NULL,
	wl_id BIGINT NOT NULL REFERENCES c2v_word_list(id),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE c2v_word_stat (
	word_id BIGINT NOT NULL REFERENCES c2v_word(id),
	user_id BIGINT NOT NULL REFERENCES c2v_user(id),
	last_training_dt TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	last_cor_ans_dt TIMESTAMP WITHOUT TIME ZONE,
	fc_next_training_dt TIMESTAMP WITHOUT TIME ZONE,
	trainings_num BIGINT NOT NULL,
	correct_answers_num BIGINT NOT NULL,
	fc_trainings_num BIGINT
);

CREATE TABLE c2v_tg_chat (
	tg_id BIGINT UNIQUE NOT NULL,
	user_id BIGINT NOT NULL REFERENCES c2v_user(id),
	state_code VARCHAR(100),
	usr_last_act_dt TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	worker VARCHAR(100),
	in_work_from TIMESTAMP WITHOUT TIME ZONE,
	wl_frgn_lang_code VARCHAR(3),
	wl_ntv_lang_code VARCHAR(3),
	wl_id BIGINT,
	word_frgn VARCHAR(2000),
	word_id BIGINT,
    words_ids VARCHAR(2000),
	exercise_code VARCHAR(100),
	trained_words_ids VARCHAR(2000),
	bot_last_msg_id BIGINT
);