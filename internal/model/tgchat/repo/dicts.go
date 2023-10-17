package repo

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tcPkg.Cmd{
	"start":        {Code: "start", DestStateCode: "main_menu"},
	"to_main_menu": {Code: "to_main_menu", DisplayLabel: "В главное меню", DestStateCode: "main_menu"},
	"create_wl":    {Code: "create_wl", DisplayLabel: "Создать список", DestStateCode: "wl_creation_frgn_lang"},
	"wl_frgn_lang": {Code: "wl_frgn_lang", DestStateCode: "wl_creation_ntv_lang"},
	"wl_ntv_lang":  {Code: "wl_ntv_lang", DestStateCode: "wl_creation_name"},
	"wl_name":      {Code: "wl_name", DestStateCode: "main_menu"},
	"all_wl":       {Code: "all_wl", DisplayLabel: "Мои списки", DestStateCode: "all_wl"},
	"wl":           {Code: "wl", DestStateCode: "main_menu"},
}

var states = map[string]*tcPkg.State{
	"start":                 {Code: "start", AvailCmds: [][]*tcPkg.Cmd{{cmds["start"]}}},
	"main_menu":             {Code: "main_menu", MsgHdr: "Главное меню", AvailCmds: [][]*tcPkg.Cmd{{cmds["create_wl"], cmds["all_wl"]}}},
	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите изучаемый язык", WaitForWLFrgnLang: true, StateCmd: cmds["wl_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите родной (базовый) язык", WaitForWLNtvLang: true, StateCmd: cmds["wl_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_name":      {Code: "wl_creation_name", MsgHdr: "Создание списка слов", MsgBody: "Введите название списка", WaitForWLName: true, NextStateCode: "main_menu", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"all_wl":                {Code: "all_wl", MsgHdr: "Мои списки", StateCmd: cmds["wl"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
}
