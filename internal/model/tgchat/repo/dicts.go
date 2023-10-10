package repo

import (
	"github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tgchat.Cmd{
	"start":    {Code: "start", DestState: states["main_menu"]},
	"to_start": {Code: "to_start", DestState: states["start"]},
}

var states = map[string]*tgchat.State{
	"start":     {Code: "start", MsgHdr: "Стартовый экран", Msg: "Доступные команды: /start", AvailCmdCodes: []string{"start"}},
	"main_menu": {Code: "main_menu", MsgHdr: "Главное меню", Msg: "Доступные команды: /to_start", AvailCmdCodes: []string{"to_start"}},
}
