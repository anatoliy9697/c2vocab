package usecases

import (
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MapToInnerUserAndSave(r res.Resources, outerU *tgbotapi.User) (u *usrPkg.User, err error) {
	u = usrPkg.MapToInner(outerU)

	var userExists bool
	if userExists, err = r.UsrRepo.IsExists(u); err == nil {
		if userExists {
			err = r.UsrRepo.Update(u)
		} else {
			err = r.UsrRepo.SaveNew(u)
		}
	}

	return u, err
}
