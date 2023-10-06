package repo

import (
	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tgchat.Cmd{
	"start": {Code: "start", DestState: states["mainmenu"]},
}

var states = map[string]*tgchat.State{
	"start":    {Code: "start", Msg: "", MsgHdr: "", MsgFtr: "", AvailCmdCodes: []string{"start"}},
	"mainmenu": {Code: "mainmenu", Msg: "Вы находитесь в главном меню", MsgHdr: "Шапка главного меню", MsgFtr: "Подвал главного меню"},
}
