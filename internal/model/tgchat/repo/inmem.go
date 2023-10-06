package repo

import (
	"errors"

	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

func (r *pgRepo) StartState() (*tgchat.State, error) {
	return states["start"], nil
}

func (r *pgRepo) StateByCode(c string) (*tgchat.State, error) {
	state, ok := states[c]
	if !ok {
		return nil, errors.New("state not found by code")
	}

	return state, nil
}
