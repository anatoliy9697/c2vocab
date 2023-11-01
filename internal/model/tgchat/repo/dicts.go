package repo

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tcPkg.Cmd{
	"start": {Code: "start", DestStateCode: "main_menu"},

	// Navigation
	"to_main_menu":  {Code: "to_main_menu", DisplayLabel: "⬅️ В главное меню", DestStateCode: "main_menu"},
	"all_wl":        {Code: "all_wl", DisplayLabel: "🗂 Мои списки", DestStateCode: "all_wl"},
	"wl":            {Code: "wl", DestStateCode: "wl"},
	"w":             {Code: "w", DestStateCode: "w"},
	"back_to_wl":    {Code: "back_to_wl", DisplayLabel: "⬅️ Назад", DestStateCode: "wl"},
	"all_w":         {Code: "all_w", DisplayLabel: "📋 Все слова списка", DestStateCode: "all_w"},
	"back_to_all_w": {Code: "back_to_all_w", DisplayLabel: "⬅️ Назад", DestStateCode: "all_w"},

	// Word list
	"create_wl":             {Code: "create_wl", DisplayLabel: "📝 Создать список", DestStateCode: "wl_creation_frgn_lang"},
	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", DestStateCode: "wl_creation_ntv_lang"},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", DestStateCode: "wl_creation_name"},
	"wl_creation_name":      {Code: "wl_creation_name", DestStateCode: "wl"},
	"edit_wl":               {Code: "edit_wl", DisplayLabel: "✏️ Редактировать", DestStateCode: "wl_edit_frgn_lang"},
	"wl_edit_frgn_lang":     {Code: "wl_edit_frgn_lang", DestStateCode: "wl_edit_ntv_lang"},
	"wl_edit_ntv_lang":      {Code: "wl_edit_ntv_lang", DestStateCode: "wl_edit_name"},
	"wl_edit_name":          {Code: "wl_edit_name", DestStateCode: "wl"},
	"delete_wl":             {Code: "delete_wl", DisplayLabel: "❌ Удалить", DestStateCode: "wl_del_confirmation"},
	"confirm_wl_del":        {Code: "confirm_wl_del", DisplayLabel: "✅ Да", DestStateCode: "all_wl"},
	"reject_wl_del":         {Code: "reject_wl_del", DisplayLabel: "❌ Нет", DestStateCode: "wl"},

	// Word
	"add_w":         {Code: "add_w", DisplayLabel: "📝 Добавить слово", DestStateCode: "w_addition_frgn"},
	"delete_w":      {Code: "delete_w", DisplayLabel: "❌ Удалить", DestStateCode: "w_del_confirmation"},
	"confirm_w_del": {Code: "confirm_w_del", DisplayLabel: "✅ Да", DestStateCode: "all_w"},
	"reject_w_del":  {Code: "reject_w_del", DisplayLabel: "❌ Нет", DestStateCode: "w"},
}

var states = map[string]*tcPkg.State{
	"start": {Code: "start", AvailCmds: [][]*tcPkg.Cmd{{cmds["start"]}}},

	// Navigation
	"main_menu": {Code: "main_menu", MsgHdr: "Главное меню", MsgBody: "Привет, {{.UsrTgFName}} {{.UsrTgLName}}!", AvailCmds: [][]*tcPkg.Cmd{{cmds["create_wl"], cmds["all_wl"]}}},
	"all_wl":    {Code: "all_wl", MsgHdr: "Мои списки", StateCmd: cmds["wl"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl":        {Code: "wl", MsgHdr: "Список слов \"{{.WLName}}\"", MsgBody: "Изучаемый язык: {{.WLFrgnLang}}\nБазовый язык: {{.WLNtvLang}}\nВсего слов: {{.WordsNum}} шт.", AvailCmds: [][]*tcPkg.Cmd{{cmds["all_w"]}, {cmds["add_w"]}, {cmds["delete_wl"], cmds["edit_wl"]}, {cmds["to_main_menu"]}}},
	"all_w":     {Code: "all_w", MsgHdr: "Слова списка \"{{.WLName}}\"", StateCmd: cmds["w"], AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w":         {Code: "w", MsgHdr: "\"{{.WordForeign}}\" - \"{{.WordNative}}\"", MsgBody: "Список слов: \"{{.WLName}}\"\nИзучаемый язык: {{.WLFrgnLang}}\nБазовый язык: {{.WLNtvLang}}", AvailCmds: [][]*tcPkg.Cmd{{cmds["delete_w"]}, {cmds["back_to_all_w"]}, {cmds["to_main_menu"]}}},

	// Word list
	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите изучаемый язык", WaitForWLFrgnLang: true, StateCmd: cmds["wl_creation_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите родной (базовый) язык", WaitForWLNtvLang: true, StateCmd: cmds["wl_creation_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_name":      {Code: "wl_creation_name", MsgHdr: "Создание списка слов", MsgBody: "Введите название списка", WaitForWLName: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_edit_frgn_lang":     {Code: "wl_edit_frgn_lang", MsgHdr: "Редактирование списка слов", MsgBody: "Выберите изучаемый язык", WaitForWLFrgnLang: true, StateCmd: cmds["wl_edit_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_edit_ntv_lang":      {Code: "wl_edit_ntv_lang", MsgHdr: "Редактирование списка слов", MsgBody: "Выберите родной (базовый) язык", WaitForWLNtvLang: true, StateCmd: cmds["wl_edit_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_edit_name":          {Code: "wl_edit_name", MsgHdr: "Редактирование списка слов", MsgBody: "Введите название списка", WaitForWLName: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_del_confirmation":   {Code: "wl_del_confirmation", MsgHdr: "Удаление списка слов", MsgBody: "Вы действительно хотите удалить список \"{{.WLName}}\"?", AvailCmds: [][]*tcPkg.Cmd{{cmds["confirm_wl_del"], cmds["reject_wl_del"]}}},

	// Word
	"w_addition_frgn":    {Code: "w_addition_frgn", MsgHdr: "Новое слово списка \"{{.WLName}}\"", MsgBody: "Введите слово на изучаемом языке ({{.WLFrgnLang}})", WaitForWFrgn: true, NextStateCode: "w_addition_ntv", AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w_addition_ntv":     {Code: "w_addition_ntv", MsgHdr: "Новое слово списка \"{{.WLName}}\"", MsgBody: "Введите перевод слова на базовом языке ({{.WLNtvLang}})", WaitForWNtv: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w_del_confirmation": {Code: "w_del_confirmation", MsgHdr: "Удаление слова", MsgBody: "Вы действительно хотите удалить слово \"{{.WordForeign}}\" - \"{{.WordNative}}\"?", AvailCmds: [][]*tcPkg.Cmd{{cmds["confirm_w_del"], cmds["reject_w_del"]}}},
}
