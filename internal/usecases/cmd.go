package usecases

import (
	"strconv"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
)

func ClearTgChaTmpFields(tc *tcPkg.Chat) {
	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil
	tc.WLId = 0
	tc.WL = nil
	tc.WordFrgn = ""
}

func SetTgChatWLFrgnLang(tc *tcPkg.Chat, langCode string) {
	tc.WLFrgnLang = commons.LangByCode(langCode)
}

func SetTgChatWLNtvLang(tc *tcPkg.Chat, langCode string) {
	tc.WLNtvLang = commons.LangByCode(langCode)
}

func CreateWL(r res.Resources, tc *tcPkg.Chat, wlName string) (err error) {
	wl := &wlPkg.WordList{
		Active:    true,
		Name:      wlName,
		FrgnLang:  tc.WLFrgnLang,
		NtvLang:   tc.WLNtvLang,
		OwnerId:   tc.UserId,
		CreatedAt: time.Now(),
	}

	if err = r.WLRepo.SaveNewWL(wl); err != nil {
		return err
	}

	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil

	return nil
}

func SetTgChatWL(r res.Resources, tc *tcPkg.Chat, wlIdStr string) (err error) {
	if tc.WLId, err = strconv.Atoi(wlIdStr); err != nil {
		return err
	}

	if tc.WL, err = r.WLRepo.WLById(tc.WLId); err != nil {
		return err
	}

	return nil
}

func DeleteWL(r res.Resources, wl *wlPkg.WordList) (err error) {
	wl.Deactivate()

	return r.WLRepo.UpdateWL(wl)
}

func SetTgChatWordFrgn(tc *tcPkg.Chat, wordFrgn string) {
	tc.WordFrgn = wordFrgn
}

func CreateWord(r res.Resources, tc *tcPkg.Chat, wordNtv string) (err error) {
	w := &wlPkg.Word{
		Foreign:   tc.WordFrgn,
		Native:    wordNtv,
		WLId:      tc.WL.Id,
		CreatedAt: time.Now(),
	}

	if err = r.WLRepo.SaveNewWord(w); err != nil {
		return err
	}

	tc.WL.WordsNum++

	tc.WordFrgn = ""

	return nil
}

func BackToWL(tc *tcPkg.Chat) {
	tc.WordFrgn = ""
}
