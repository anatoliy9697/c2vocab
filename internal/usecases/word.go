package usecases

import (
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
)

func RegistrateWordTraining(r res.Resources, w *wlPkg.Word, userId int, isAnswerCorrect bool) (err error) {
	isStatExists := false

	if isStatExists, err = r.WLRepo.IsWordStatExists(w.Id, userId); err != nil {
		return err
	}
	if isStatExists {
		err = r.WLRepo.RegistrateWordTraining(w.Id, userId, isAnswerCorrect, w.MemPercentageDowngrade())
	} else {
		err = r.WLRepo.CreateWordStat(w.Id, userId, isAnswerCorrect)
	}

	return err
}
