package usecases

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	"github.com/jackc/pgx/v5"
)

func ClearTgChaTmpFields(tc *tcPkg.Chat) {
	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil
	ClearWLFields(tc)
	ClearWordCreationFields(tc)
	ClearWordFields(tc)
	ClearExerciseFields(tc)
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

	tc.WLId = wl.Id
	tc.WL = wl

	tc.WLFrgnLang = nil
	tc.WLNtvLang = nil

	return nil
}

func EditWL(r res.Resources, tc *tcPkg.Chat, wlName string) (err error) {
	tc.WL.Name = wlName
	tc.WL.FrgnLang = tc.WLFrgnLang
	tc.WL.NtvLang = tc.WLNtvLang

	if err = r.WLRepo.UpdateWL(tc.WL); err != nil {
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

func ClearWordCreationFields(tc *tcPkg.Chat) {
	tc.WordFrgn = ""
}

func SetTgChatWord(r res.Resources, tc *tcPkg.Chat, wordIdStr string) (err error) {
	if tc.WordId, err = strconv.Atoi(wordIdStr); err != nil {
		return err
	}

	if tc.Word, err = r.WLRepo.WordById(tc.WordId); err != nil {
		return err
	}

	return nil
}

func ClearWordFields(tc *tcPkg.Chat) {
	tc.WordId = 0
	tc.Word = nil
}

func DeleteWord(r res.Resources, w *wlPkg.Word) (err error) {
	w.Deactivate()

	return r.WLRepo.UpdateWord(w)
}

func SetTgChatExercise(r res.Resources, tc *tcPkg.Chat, excersiceCode string) (err error) {
	tc.ExcersiceCode = excersiceCode

	if tc.Excersice, err = r.TcRepo.ExcersiceByCode(excersiceCode); err != nil {
		return err
	}

	if tc.Word, err = r.WLRepo.RandActiveWordByWLIdAndExcludedIds(tc.WL.Id, ""); err != nil {
		return err
	}
	tc.WordId = tc.Word.Id

	return nil
}

func ClearExerciseFields(tc *tcPkg.Chat) {
	tc.ExcersiceCode = ""
	tc.Excersice = nil
	tc.WordId = 0
	tc.Word = nil
	tc.TrainedWordsIds = ""
}

func ClearWLFields(tc *tcPkg.Chat) {
	tc.WLId = 0
	tc.WL = nil
}

func ProcessUserTaskDataInput(r res.Resources, tc *tcPkg.Chat, usrAnswer string) (err error) {
	usrAnswer = strings.ToLower(usrAnswer)

	prevTaskResult := "Правильно!"

	switch tc.ExcersiceCode {

	case "write_frgn":
		if usrAnswer != strings.ToLower(tc.Word.Foreign) {
			prevTaskResult = "Неправильно. Правильный ответ: \"" + tc.Word.Foreign + "\""
		}

	default:

	}

	if tc.TrainedWordsIds == "" {
		tc.TrainedWordsIds = strconv.Itoa(tc.WordId)
	} else {
		tc.TrainedWordsIds += ", " + strconv.Itoa(tc.WordId)
	}
	tc.PrevTaskResult = prevTaskResult
	if tc.Word, err = r.WLRepo.RandActiveWordByWLIdAndExcludedIds(tc.WL.Id, tc.TrainedWordsIds); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // TODO: тут лучше вынести в отдельную функцию
			var s *tcPkg.State
			if s, err = r.TcRepo.StateByCode("xrcs_finish"); err != nil {
				return err
			}
			tc.SetState(s)
			ClearExerciseFields(tc)
		} else {
			return err
		}
	}
	if tc.Word != nil {
		tc.WordId = tc.Word.Id
	}

	return nil
}
