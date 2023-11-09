package usecases

import (
	"strconv"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func WLFrgnLangTgInlineKeyboard(tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton) {
	rowsCount := len(commons.AvailLangs)
	ikRows = make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i, lang := range commons.AvailLangs {
		ikRows[i] = make([]tgbotapi.InlineKeyboardButton, 1)
		ikRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.Cmd.Code+" "+lang.Code)
	}

	return ikRows
}

func WLNtvLangTgInlineKeyboard(tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton) {
	rowsCount := len(commons.AvailLangs) - 1
	ikRows = make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	i := 0
	for _, lang := range commons.AvailLangs {
		if lang.Code != tc.WLFrgnLang.Code {
			ikRows[i] = make([]tgbotapi.InlineKeyboardButton, 1)
			ikRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.Cmd.Code+" "+lang.Code)
			i++
		}
	}

	return ikRows
}

func AllWLTgInlineKeyboard(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	var wls []*wlPkg.WordList
	if wls, err = r.WLRepo.ActiveWLByOwnerId(tc.UserId); err != nil {
		return nil, err
	}

	var ikRow []tgbotapi.InlineKeyboardButton
	i := 1
	for _, wl := range wls {
		ikRow = []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i)+". \""+wl.Name+"\"", tc.State.Cmd.Code+" "+strconv.Itoa(int(wl.Id)))}
		ikRows = append(ikRows, ikRow)
		i++
	}

	return ikRows, nil
}

func AllWordsTgInlineKeyboard(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	var words []*wlPkg.Word
	if words, err = r.WLRepo.ActiveWordsByWLId(tc.WL.Id); err != nil {
		return nil, err
	}

	var ikRow []tgbotapi.InlineKeyboardButton
	i := 1
	for _, word := range words {
		ikRow = []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i)+". \""+word.Foreign+"\" - \""+word.Native+"\"", tc.State.Cmd.Code+" "+strconv.Itoa(int(word.Id)))}
		ikRows = append(ikRows, ikRow)
		i++
	}

	return ikRows, nil
}

func AllExercisesTgInlineKeyboard(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton) {
	xrcses := r.TcRepo.AllExercises()

	var ikRow []tgbotapi.InlineKeyboardButton
	i := 1
	for _, xrcs := range xrcses {
		ikRow = []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i)+". "+xrcs.Name, tc.State.Cmd.Code+" "+xrcs.Code)}
		ikRows = append(ikRows, ikRow)
		i++
	}

	return ikRows
}

func WordSelectionExerciseTgInlineKeyboard(r res.Resources, tc *tcPkg.Chat, frgnOpts bool) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	var ansOpts []wlPkg.AnswerOption

	if frgnOpts {
		ansOpts, err = r.WLRepo.WordSelectionAnswerOptions(tc.Word, true, tc.WL.FrgnLang.Code, tc.UserId, 4)
	} else {
		ansOpts, err = r.WLRepo.WordSelectionAnswerOptions(tc.Word, false, tc.WL.NtvLang.Code, tc.UserId, 4)
	}
	if err != nil {
		return nil, err
	}

	ikRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)
	for _, ansOpt := range ansOpts {
		ikRow = append(ikRow, tgbotapi.NewInlineKeyboardButtonData(ansOpt.Answer, "ans "+ansOpt.IsCorrect))
		if len(ikRow) == 2 {
			ikRows = append(ikRows, ikRow)
			ikRow = make([]tgbotapi.InlineKeyboardButton, 0, 2)
		}
	}
	if len(ikRow) == 2 {
		ikRows = append(ikRows, ikRow)
	}

	return ikRows, nil
}

func TgInlineKeyboradExerciseCmdRows(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	switch tc.Excersice.Code {
	case "select_frgn":
		if ikRows, err = WordSelectionExerciseTgInlineKeyboard(r, tc, true); err != nil {
			return nil, err
		}
	case "select_ntv":
		if ikRows, err = WordSelectionExerciseTgInlineKeyboard(r, tc, false); err != nil {
			return nil, err
		}
	default:
	}

	return ikRows, nil
}

func TgInlineKeyboradStateCmdRows(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	switch tc.State.Cmd.Code {
	case "wl_creation_frgn_lang", "wl_editing_frgn_lang":
		ikRows = WLFrgnLangTgInlineKeyboard(tc)
	case "wl_creation_ntv_lang", "wl_editing_ntv_lang":
		ikRows = WLNtvLangTgInlineKeyboard(tc)
	case "wl":
		if ikRows, err = AllWLTgInlineKeyboard(r, tc); err != nil {
			return nil, err
		}
	case "w":
		if ikRows, err = AllWordsTgInlineKeyboard(r, tc); err != nil {
			return nil, err
		}
	case "xrcs":
		ikRows = AllExercisesTgInlineKeyboard(r, tc)
	default:
	}

	return ikRows, nil
}

func TgInlineKeyboradAvailCmdsRows(r res.Resources, tc *tcPkg.Chat) [][]tgbotapi.InlineKeyboardButton {
	emptyWL := (tc.WL != nil && tc.WL.WordsNum == 0)

	availCmds := tc.State.AvailCmdsByFlgs(emptyWL)

	inlineKeyboardRows := make([][]tgbotapi.InlineKeyboardButton, len(availCmds))
	for i, cmdsRow := range availCmds {
		inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, len(cmdsRow))
		for j, cmd := range cmdsRow {
			inlineKeyboardRows[i][j] = cmd.TgButton()
		}
	}

	return inlineKeyboardRows
}

func TgInlineKeyboradMarkup(r res.Resources, tc *tcPkg.Chat) (ik tgbotapi.InlineKeyboardMarkup, err error) {
	ikRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	if tc.Excersice != nil && tc.Excersice.Cmd != nil && tc.State.Code == "xrcs" {
		ikExerciseCmdRows, err := TgInlineKeyboradExerciseCmdRows(r, tc)
		if err != nil {
			return ik, err
		}
		if len(ikExerciseCmdRows) > 0 {
			ikRows = append(ikRows, ikExerciseCmdRows...)
		}
	}

	if tc.State.Cmd != nil {
		ikStateCmdRows, err := TgInlineKeyboradStateCmdRows(r, tc)
		if err != nil {
			return ik, err
		}
		if len(ikStateCmdRows) > 0 {
			ikRows = append(ikRows, ikStateCmdRows...)
		}
	}

	ikAvailCmdsRows := TgInlineKeyboradAvailCmdsRows(r, tc)
	if len(ikAvailCmdsRows) > 0 {
		ikRows = append(ikRows, ikAvailCmdsRows...)
	}

	if len(ikRows) > 0 {
		ik = tgbotapi.NewInlineKeyboardMarkup(ikRows...)
	}

	return ik, nil
}
