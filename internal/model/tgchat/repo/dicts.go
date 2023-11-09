package repo

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

var cmds = map[string]*tcPkg.Cmd{
	// Navigation
	"start":          {Code: "start", DestStateCode: "main_menu"},
	"to_main_menu":   {Code: "to_main_menu", DisplayLabel: "⬅️ В главное меню", DestStateCode: "main_menu"},
	"back_to_all_wl": {Code: "back_to_all_wl", DisplayLabel: "⬅️ Назад", DestStateCode: "all_wl"},
	"back_to_wl":     {Code: "back_to_wl", DisplayLabel: "⬅️ Назад", DestStateCode: "wl"},
	"back_to_all_w":  {Code: "back_to_all_w", DisplayLabel: "⬅️ Назад", DestStateCode: "all_w"},
	"finish_xrcs":    {Code: "finish_xrcs", DisplayLabel: "🏁 Закончить", DestStateCode: "wl"},

	// Word list
	"wl":                    {Code: "wl", DestStateCode: "wl"},
	"all_wl":                {Code: "all_wl", DisplayLabel: "🗂 Мои списки", DestStateCode: "all_wl"},
	"create_wl":             {Code: "create_wl", DisplayLabel: "📝 Создать список", DestStateCode: "wl_creation_frgn_lang"},
	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", DestStateCode: "wl_creation_ntv_lang"},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", DestStateCode: "wl_creation_name"},
	"wl_creation_name":      {Code: "wl_creation_name", DestStateCode: "wl"},
	"edit_wl":               {Code: "edit_wl", DisplayLabel: "✏️ Редактировать", DestStateCode: "wl_editing_frgn_lang"},
	"wl_editing_frgn_lang":  {Code: "wl_editing_frgn_lang", DestStateCode: "wl_editing_ntv_lang"},
	"wl_editing_ntv_lang":   {Code: "wl_editing_ntv_lang", DestStateCode: "wl_editing_name"},
	"wl_editing_name":       {Code: "wl_editing_name", DestStateCode: "wl"},
	"delete_wl":             {Code: "delete_wl", DisplayLabel: "❌ Удалить", DestStateCode: "wl_del_confirmation"},
	"confirm_wl_del":        {Code: "confirm_wl_del", DisplayLabel: "✅ Да", DestStateCode: "all_wl"},
	"reject_wl_del":         {Code: "reject_wl_del", DisplayLabel: "❌ Нет", DestStateCode: "wl"},

	// Word
	"w":             {Code: "w", DestStateCode: "w"},
	"all_w":         {Code: "all_w", DisplayLabel: "📋 Все слова списка", DestStateCode: "all_w", NotEmptyWLOnly: true},
	"add_w":         {Code: "add_w", DisplayLabel: "📝 Добавить слово", DestStateCode: "w_addition_frgn"},
	"delete_w":      {Code: "delete_w", DisplayLabel: "❌ Удалить", DestStateCode: "w_del_confirmation"},
	"confirm_w_del": {Code: "confirm_w_del", DisplayLabel: "✅ Да", DestStateCode: "all_w"},
	"reject_w_del":  {Code: "reject_w_del", DisplayLabel: "❌ Нет", DestStateCode: "w"},

	// Learning
	"learn_wl": {Code: "learn_wl", DisplayLabel: "🧠 Учить", DestStateCode: "all_exercises", NotEmptyWLOnly: true},
	"xrcs":     {Code: "xrcs", DestStateCode: "xrcs"},
	"ans":      {Code: "ans"},
}

var states = map[string]*tcPkg.State{
	// Navigation
	"start":     {Code: "start", AvailCmds: [][]*tcPkg.Cmd{{cmds["start"]}}},
	"main_menu": {Code: "main_menu", MsgHdr: "Главное меню", MsgBody: "Привет, {{.UsrTgFName}} {{.UsrTgLName}}!", AvailCmds: [][]*tcPkg.Cmd{{cmds["create_wl"], cmds["all_wl"]}}},

	// Word list
	"wl":                    {Code: "wl", MsgHdr: "Список слов \"{{.WLName}}\"", MsgBody: "Изучаемый язык: {{.WLFrgnLang}}\nБазовый язык: {{.WLNtvLang}}\nВсего слов: {{.WordsNum}} шт.", AvailCmds: [][]*tcPkg.Cmd{{cmds["learn_wl"]}, {cmds["all_w"]}, {cmds["add_w"]}, {cmds["delete_wl"], cmds["edit_wl"]}, {cmds["back_to_all_wl"]}, {cmds["to_main_menu"]}}},
	"all_wl":                {Code: "all_wl", MsgHdr: "Мои списки", Cmd: cmds["wl"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_frgn_lang": {Code: "wl_creation_frgn_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите изучаемый язык", Cmd: cmds["wl_creation_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_ntv_lang":  {Code: "wl_creation_ntv_lang", MsgHdr: "Создание списка слов", MsgBody: "Выберите базовый (родной) язык", Cmd: cmds["wl_creation_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_creation_name":      {Code: "wl_creation_name", MsgHdr: "Создание списка слов", MsgBody: "Введите название списка", WaitForDataInput: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_editing_frgn_lang":  {Code: "wl_editing_frgn_lang", MsgHdr: "Редактирование списка слов", MsgBody: "Выберите изучаемый язык", Cmd: cmds["wl_editing_frgn_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_editing_ntv_lang":   {Code: "wl_editing_ntv_lang", MsgHdr: "Редактирование списка слов", MsgBody: "Выберите базовый (родной) язык", Cmd: cmds["wl_editing_ntv_lang"], AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_editing_name":       {Code: "wl_editing_name", MsgHdr: "Редактирование списка слов", MsgBody: "Введите название списка", WaitForDataInput: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["to_main_menu"]}}},
	"wl_del_confirmation":   {Code: "wl_del_confirmation", MsgHdr: "Удаление списка слов", MsgBody: "Вы действительно хотите удалить список \"{{.WLName}}\"?", AvailCmds: [][]*tcPkg.Cmd{{cmds["confirm_wl_del"], cmds["reject_wl_del"]}}},

	// Word
	"w":                  {Code: "w", MsgHdr: "\"{{.WordForeign}}\" - \"{{.WordNative}}\"", MsgBody: "Список слов: \"{{.WLName}}\"\nИзучаемый язык: {{.WLFrgnLang}}\nБазовый язык: {{.WLNtvLang}}", AvailCmds: [][]*tcPkg.Cmd{{cmds["delete_w"]}, {cmds["back_to_all_w"]}, {cmds["to_main_menu"]}}},
	"all_w":              {Code: "all_w", MsgHdr: "Слова списка \"{{.WLName}}\"", Cmd: cmds["w"], AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w_addition_frgn":    {Code: "w_addition_frgn", MsgHdr: "Новое слово списка \"{{.WLName}}\"", MsgBody: "Введите слово на изучаемом языке ({{.WLFrgnLang}})", WaitForDataInput: true, NextStateCode: "w_addition_ntv", AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w_addition_ntv":     {Code: "w_addition_ntv", MsgHdr: "Новое слово списка \"{{.WLName}}\"", MsgBody: "Введите перевод слова на базовом языке ({{.WLNtvLang}})", WaitForDataInput: true, NextStateCode: "wl", AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"w_del_confirmation": {Code: "w_del_confirmation", MsgHdr: "Удаление слова", MsgBody: "Вы действительно хотите удалить слово \"{{.WordForeign}}\" - \"{{.WordNative}}\"?", AvailCmds: [][]*tcPkg.Cmd{{cmds["confirm_w_del"], cmds["reject_w_del"]}}},

	// Learning
	"all_exercises": {Code: "all_exercises", MsgHdr: "Изучение списка слов \"{{.WLName}}\"", MsgBody: "Выберите упражнение", Cmd: cmds["xrcs"], AvailCmds: [][]*tcPkg.Cmd{{cmds["back_to_wl"]}, {cmds["to_main_menu"]}}},
	"xrcs":          {Code: "xrcs", MsgBody: "{{.ExerciseTaskText}}", AvailCmds: [][]*tcPkg.Cmd{{cmds["finish_xrcs"]}}},
	"xrcs_finish":   {Code: "xrcs_finish", MsgBody: "{{.PrevTaskResult}}На этом пока все! :)", AvailCmds: [][]*tcPkg.Cmd{{cmds["finish_xrcs"]}, {cmds["to_main_menu"]}}},
}

var exercises = map[string]*tcPkg.Excersice{
	"write_frgn":  {Code: "write_frgn", Name: "Ввод слова", TaskText: "{{.PrevTaskResult}}Введите слово \"{{.WordNative}}\" на изучаемом ({{.WLFrgnLang}}) языке", WaitForDataInput: true},
	"select_frgn": {Code: "select_frgn", Name: "Выбор слова", TaskText: "{{.PrevTaskResult}}Выберите слово \"{{.WordNative}}\" на изучаемом ({{.WLFrgnLang}}) языке", Cmd: cmds["ans"]},
	"select_ntv":  {Code: "select_ntv", Name: "Выбор перевода", TaskText: "{{.PrevTaskResult}}Выберите перевод слова \"{{.WordForeign}}\" на базовом ({{.WLNtvLang}}) языке", Cmd: cmds["ans"]},
}
