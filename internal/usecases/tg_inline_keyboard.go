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
		ikRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.StateCmd.Code+" "+lang.Code)
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
			ikRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.StateCmd.Code+" "+lang.Code)
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
	for _, wl := range wls {
		ikRow = []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(wl.Name, tc.State.StateCmd.Code+" "+strconv.Itoa(int(wl.Id)))}
		ikRows = append(ikRows, ikRow)
	}

	return ikRows, nil
}

func TgInlineKeyboradStateCmdRows(r res.Resources, tc *tcPkg.Chat) (ikRows [][]tgbotapi.InlineKeyboardButton, err error) {
	switch {
	case tc.State.WaitForWLFrgnLang:
		ikRows = WLFrgnLangTgInlineKeyboard(tc)
	case tc.State.WaitForWLNtvLang:
		ikRows = WLNtvLangTgInlineKeyboard(tc)
	case tc.State.StateCmd != nil && tc.State.StateCmd.Code == "wl":
		if ikRows, err = AllWLTgInlineKeyboard(r, tc); err != nil {
			return nil, err
		}
	}

	return ikRows, nil
}

func TgInlineKeyboradMarkup(r res.Resources, tc *tcPkg.Chat) (ik tgbotapi.InlineKeyboardMarkup, err error) {
	ikRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	ikStateCmdRows, err := TgInlineKeyboradStateCmdRows(r, tc)
	if err != nil {
		return ik, err
	}
	if len(ikStateCmdRows) > 0 {
		ikRows = append(ikRows, ikStateCmdRows...)
	}

	ikAvailCmdsRows := tc.State.TgInlineKeyboradAvailCmdsRows()
	if len(ikAvailCmdsRows) > 0 {
		ikRows = append(ikRows, ikAvailCmdsRows...)
	}

	if len(ikRows) > 0 {
		ik = tgbotapi.NewInlineKeyboardMarkup(ikRows...)
	}

	return ik, nil
}
