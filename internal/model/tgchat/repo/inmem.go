package repo

import (
	"errors"
	"text/template"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

func initStateMsgTmpls() (err error) {
	var tmplContent string
	var tmpl *template.Template

	for _, s := range states {
		tmplContent = s.OutMsgTmplContent()
		if tmpl, err = template.New(s.Code).Parse(tmplContent); err != nil {
			return err
		}
		s.MsgTmpl = tmpl
	}

	return nil
}

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
