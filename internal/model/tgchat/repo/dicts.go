package repo

import (
	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tgchat.Cmd{
	"start": {Code: "start", DestState: states["main_menu"]},
}

var states = map[string]*tgchat.State{
	"start":     {Code: "start", Msg: "", MsgHdr: "", MsgFtr: "", AvailCmdCodes: []string{"start"}},
	"main_menu": {Code: "main_menu", Msg: "Вы находитесь в главном меню", MsgHdr: "Шапка главного меню", MsgFtr: "Подвал главного меню"},
}
