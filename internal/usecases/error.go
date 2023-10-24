package usecases

import (
	"errors"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
)

func ProcessErr(r res.Resources, tc *tcPkg.Chat, usr *usrPkg.User, err error) {
	if err == nil {
		return
	}

	var validationErr *tcPkg.IncMsgValidationErr
	errText := ""
	if errors.As(err, &validationErr) {
		errText = "Пользовательская ошибка: " + validationErr.Error()
	} else {
		r.Logger.Error(err.Error())
		errText = "Возникла системная ошибка. Попробуйте позднее"
	}

	if tc != nil {
		prevBotLastMsgId := tc.BotLastMsgId

		// Getting tgChat without any changes made in RAM
		tc, _ = r.TcRepo.TgChatByUserId(tc.UserId)
		if tc != nil {
			tc.User = usr

			if tc.WLId != 0 {
				tc.WL, _ = r.WLRepo.WLById(tc.WLId)
			}

			// Sending replay message, got by non-changed by current request, tgChat state with error info
			tc.BotLastMsgId = prevBotLastMsgId
			SendReplyMsg(r, tc, errText)
			r.TcRepo.UpdateTgChat(tc)
		}
	}
}
