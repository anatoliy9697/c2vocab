package usecases

import (
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	usrRepo "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MapToInnerUserAndSave(r usrRepo.Repo, outerU *tgbotapi.User) (u *usrPkg.User, err error) {
	u = usrPkg.MapToInner(outerU)

	var userExists bool
	if userExists, err = r.IsExists(u); err == nil {
		if userExists {
			err = r.Update(u)
		} else {
			err = r.SaveNew(u)
		}
	}

	return u, err
}
