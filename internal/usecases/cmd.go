package usecases

import (
	"strconv"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
)

func ClearTgChaTmpFields(tc *tcPkg.Chat) {
	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil
	tc.WLId = 0
	tc.WL = nil
}

func SetTgChatWLFrgnLang(tc *tcPkg.Chat, langCode string) {
	tc.WLFrgnLang = commons.LangByCode(langCode)
}

func SetTgChatWLNtvLang(tc *tcPkg.Chat, langCode string) {
	tc.WLNtvLang = commons.LangByCode(langCode)
}

func CreateWL(wlR wlRepo.Repo, tc *tcPkg.Chat, wlName string) (err error) {
	wl := &wlPkg.WordList{
		Active:    true,
		Name:      wlName,
		FrgnLang:  tc.WLFrgnLang,
		NtvLang:   tc.WLNtvLang,
		OwnerId:   tc.UserId,
		CreatedAt: time.Now(),
	}

	if err = wlR.SaveNewWL(wl); err != nil {
		return err
	}

	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil

	return nil
}

func SetTgChatWL(wlR wlRepo.Repo, tc *tcPkg.Chat, wlIdStr string) (err error) {
	var id int
	if id, err = strconv.Atoi(wlIdStr); err != nil {
		return err
	}

	tc.WLId = int32(id)

	if tc.WL, err = wlR.WLById(tc.WLId); err != nil {
		return err
	}

	return nil
}
