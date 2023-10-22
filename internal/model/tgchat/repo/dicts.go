package repo

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tcPkg.Cmd{
	"start":          {Code: "start", DestStateCode: "main_menu"},
	"to_main_menu":   {Code: "to_main_menu", DisplayLabel: "🔙 В главное меню", DestStateCode: "main_menu"},
	"create_wl":      {Code: "create_wl", DisplayLabel: "📝 Создать список", DestStateCode: "wl_creation_frgn_lang"},
	"wl_frgn_lang":   {Code: "wl_frgn_lang", DestStateCode: "wl_creation_ntv_lang"},
	"wl_ntv_lang":    {Code: "wl_ntv_lang", DestStateCode: "wl_creation_name"},
	"wl_name":        {Code: "wl_name", DestStateCode: "main_menu"},
	"all_wl":         {Code: "all_wl", DisplayLabel: "📜 Мои списки", DestStateCode: "all_wl"},
	"wl":             {Code: "wl", DestStateCode: "wl"},
	"delete_wl":      {Code: "delete_wl", DisplayLabel: "Удалить список", DestStateCode: "wl_del_confirmation"}, // TODO: добавить иконку
	"confirm_wl_del": {Code: "confirm_wl_del", DisplayLabel: "✅ Да", DestStateCode: "all_wl"},
	"reject_wl_del":  {Code: "reject_wl_del", DisplayLabel: "❌ Нет", DestStateCode: "wl"},
}

var states = map[string]*tcPkg.State{
	"start": {Code: "start", AvailCmds: [][]*tcPkg.Cmd{{cmds["start"]}}},

	"main_menu": {Code: "main_menu", MsgHdr: "Главное меню", AvailCmds: [][]*tcPkg.Cmd{{cmds["create_wl"], cmds["all_wl"]}}},

	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите изучаемый язык", WaitForWLFrgnLang: true, StateCmd: cmds["wl_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите родной (базовый) язык", WaitForWLNtvLang: true, StateCmd: cmds["wl_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_name":      {Code: "wl_creation_name", MsgHdr: "Создание списка слов", MsgBody: "Введите название списка", WaitForWLName: true, NextStateCode: "main_menu", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},

	"wl": {Code: "wl", MsgHdr: "Список слов \"{{.WLName}}\"", MsgBody: "Всего слов: 0 шт.", AvailCmds: [][]*tcPkg.Cmd{{cmds["delete_wl"]}, {cmds["to_main_menu"]}}},

	"all_wl": {Code: "all_wl", MsgHdr: "Мои списки", StateCmd: cmds["wl"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},

	"wl_del_confirmation": {Code: "wl_del_confirmation", MsgHdr: "Удаление списка слов", MsgBody: "Вы действительно хотите удалить список \"{{.WLName}}\"?", AvailCmds: [][]*tcPkg.Cmd{{cmds["confirm_wl_del"], cmds["reject_wl_del"]}}},
}
