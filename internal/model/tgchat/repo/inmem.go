package repo

import (
	"errors"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

func (r pgRepo) StartState() (*tcPkg.State, error) {
	return states["start"], nil
}

func (r pgRepo) StateByCode(c string) (*tcPkg.State, error) {
	state, ok := states[c]
	if !ok {
		return nil, errors.New("state not found by code")
	}

	return state, nil
}

func (r pgRepo) CmdByCode(c string) (*tcPkg.Cmd, error) {
	cmd, ok := cmds[c]
	if !ok {
		return nil, errors.New("command not found by code")
	}

	return cmd, nil
}
